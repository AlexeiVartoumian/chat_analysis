package handlers

import (
	"api/auth"
	"api/models"
	"api/repository/sqlconnect"
	"api/utils"
	"bufio"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func updateLambda(search_term string) error {

	cmd := exec.Command("/venv/bin/python3", "/home/ubuntu/update_search.py", search_term)

	if err := cmd.Run(); err != nil {
		return utils.ErrorHandler(err, "program did not execute")
	}
	return nil

}

func SemanticSearch(w http.ResponseWriter, r *http.Request) {

	//
	if r.Method != http.MethodPost {
		http.Error(w, r.Method, http.StatusBadRequest)
		http.Error(w, "mthod not allowed", http.StatusMethodNotAllowed)
	}

	//todo sanitize input
	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Invalid requestr body", http.StatusBadRequest)
	}
	defer r.Body.Close()

	search_term := string(body)

	err = updateLambda(search_term)

	if err != nil {
		log.Println(err)
		http.Error(w, "Problem executing", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Status string `json:"status"`
	}{
		Status: fmt.Sprintf(" successfully sent request for %s", search_term),
	}
	json.NewEncoder(w).Encode(response)

}
func CompanyUrlOnly(w http.ResponseWriter, r *http.Request) {

	CompanyLinks, err := sqlconnect.OnlyUrlOnCompanySite()

	if err != nil {
		http.Error(w, "Invalid request Body", http.StatusBadRequest)
		return
	}
	fmt.Println("here is yout data", CompanyLinks)

	response := struct {
		Status string                        `json:"status"`
		Count  int                           `json:"count"`
		Data   []sqlconnect.UrlOnCompanySite `json:"data"`
	}{
		Status: "success",
		Count:  len(CompanyLinks),
		Data:   CompanyLinks,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// TODO make middleware function so not anyone can do
func SeekExpiredRoles(w http.ResponseWriter, r *http.Request) {

	LiveRoles, err := sqlconnect.SeekExpired()

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	fmt.Println("here is your data", LiveRoles)

	payload, _ := json.Marshal(LiveRoles)
	auto := "false"
	cmd := exec.Command("python3", "/home/ubuntu/backfill.py", "live", auto)

	cmd.Stdin = bytes.NewReader(payload)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("backfill.py failed %v", err)
	}

}

func SeekReopenedRoles(w http.ResponseWriter, r *http.Request) {

	LiveRoles, err := sqlconnect.SeekReopened()

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	fmt.Println("here is your data", LiveRoles)

	payload, _ := json.Marshal(LiveRoles)
	auto := "false"
	cmd := exec.Command("python3", "/home/ubuntu/backfill.py", "suspended", auto)

	cmd.Stdin = bytes.NewReader(payload)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("backfill.py failed %v", err)
	}

}

func GetLastThreeDays(w http.ResponseWriter, r *http.Request) {
	recentJobs, err := sqlconnect.LastThreeDaysJobs()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("here is yout data", recentJobs)
	response := struct {
		Status string                     `json:"status"`
		Count  int                        `json:"count"`
		Data   []sqlconnect.LastThreeDays `json:"data"`
	}{
		Status: "success",
		Count:  len(recentJobs),
		Data:   recentJobs,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func HandleQuery(w http.ResponseWriter, r *http.Request) {

	var req struct {
		Query string        `json:"query"`
		Args  []interface{} `json:"args"`
	}
	decoder := json.NewDecoder(r.Body)
	decoder.UseNumber()

	if err := decoder.Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	results, cols, err := sqlconnect.QueryToJson(req.Query, req.Args...)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sqlconnect.QueryResponse{
		Columns: cols,
		Count:   len(results),
		Rows:    results,
	})
}

func PostApiKey(w http.ResponseWriter, r *http.Request) {

	generator := auth.NewAPIKeyGenerator()
	hasher := auth.NewKeyHasher()

	db, err := sqlconnect.ConnectDb()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store := sqlconnect.NewPostgresStore(db)
	// Generate new key
	fullKey, keyID, err := generator.Generate()
	if err != nil {
		http.Error(w, "Failed to generate key", http.StatusInternalServerError)
		return
	}

	lordOfTheRings, sqlerrnorows := store.ThereCanBeOnlyOne()

	if lordOfTheRings == -1 {
		utils.ErrorHandler(sqlerrnorows, "arrg this function did not fully execute")
		http.Error(w, "apologies this did not work ", http.StatusBadRequest)

		return

	}
	// if lordOfTheRings != sql.ErrNoRows {
	// 	utils.ErrorHandler(lordOfTheRings, "arrg")
	// 	http.Error(w, "welcome to db dear guest :) here is your api key", http.StatusBadRequest)
	// 	// do something
	// 	return
	// }
	// Hash the key for storage
	hashedKey, err := hasher.Hash(fullKey)
	if err != nil {
		http.Error(w, "Failed to hash key", http.StatusInternalServerError)
		return
	}
	//TODO fix userid creation logic right now hardcoding
	// Create key record
	var scopecheck models.Scope
	var UserId string
	scopecheck = models.ScopeRead

	if sqlerrnorows == sql.ErrNoRows {
		scopecheck = models.ScopeAdmin
		UserId = "00000000-0000-0000-0000-000000000001"
	} else {

		UserId = fmt.Sprintf("00000000-0000-0000-0000-00000000000%s", strconv.Itoa(lordOfTheRings+1))
	}

	//TODO refactor! delete me endpoint
	if lordOfTheRings > 2 {

		http.Error(w, "maximum db users reached ! plz use this endpoint to delete your key. /deleteme", http.StatusBadRequest)
		// do something
		return
	}

	apiKey := &models.APIKey{
		KeyID:     keyID,
		HashedKey: hashedKey,
		//Name:      r.FormValue("name"),
		Name:   fmt.Sprintf("mykey_%s", strconv.Itoa(1+lordOfTheRings)),
		UserID: UserId, // temp placeholder
		//ProjectID: r.FormValue("project_id"),
		ProjectID: "00000000-0000-0000-0000-000000000002",
		Scopes:    []models.Scope{scopecheck},
		RateLimit: 1000,
		IsActive:  true,
	}

	if err := store.Create(r.Context(), apiKey); err != nil {
		log.Println("Failed to store key:", err)
		http.Error(w, "Failed to store key", http.StatusInternalServerError)
		return
	}

	// Return the full key to the user (only shown once!)
	w.Header().Set("Content-Type", "application/json")
	if sqlerrnorows == sql.ErrNoRows {

		w.Write([]byte(`{ hey you admin "api_key": "` + fullKey + `", "key_id": "` + keyID + `"}`))
	} else {

		w.Write([]byte(`{welcome to db dear guest :) here is your api key : "api_key": "` + fullKey + `", "key_id": "` + keyID + `"}`))
	}

}

func InsertToDb(w http.ResponseWriter, r *http.Request) {

	SeenFileKeys, err := sqlconnect.GetKeys()

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	fmt.Println("here is your data", len(SeenFileKeys))

	file, err := os.Create("/home/ubuntu/keystest.txt")

	if err != nil {
		http.Error(w, "Internal server error writing file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, line := range SeenFileKeys {
		fmt.Fprintln(writer, line)
	}

	if err := writer.Flush(); err != nil {
		http.Error(w, "internal server error flushing file", http.StatusInternalServerError)
		return
	}

}

func SqsBlaster(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, r.Method, http.StatusBadRequest)
		http.Error(w, "method not allowed ", http.StatusMethodNotAllowed)
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Invalid requestr body", http.StatusBadRequest)
	}
	defer r.Body.Close()

	first_run_and_number_accounts := string(body)
	first_run := false
	number_accounts := 0
	//TODO dont use data for identification use something else
	if strings.Contains(first_run_and_number_accounts, "first") {

		first_run = true
		number_accounts, err = strconv.Atoi(strings.Split(first_run_and_number_accounts, " ")[1])

		if err != nil {
			log.Println(err)
			http.Error(w, "Problem executing most likely data in unexpected format", http.StatusInternalServerError)
		}

	} else {
		number_accounts, err = strconv.Atoi(first_run_and_number_accounts)
		if err != nil {
			log.Println(err)
			http.Error(w, "Problem executing most likely data in unexpected format", http.StatusInternalServerError)
		}
	}

	fmt.Println(first_run_and_number_accounts)

	SearchTerms, err := sqlconnect.GetSearchTerms(first_run, number_accounts, "SEARCH_TERM")

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if SearchTerms == nil {
		fmt.Println("Job Done either second last run or last run ")
		return
	}

	fmt.Println("here is your data", SearchTerms, len(SearchTerms))

	payload := SendToSqs(SearchTerms)

	w.Header().Set("Content-Type", "application/json")

	response := struct {
		Status     string `json:"status"`
		StatusCode int
	}{
		Status:     fmt.Sprintf(" Successfully blasted sqs. have a good day %d", len(payload)),
		StatusCode: 200,
	}

	json.NewEncoder(w).Encode(response)

}

func SendToSqs(SearchTerms []sqlconnect.SearchTerm) []byte {
	payload, _ := json.Marshal(SearchTerms)

	cmd := exec.Command("python3", "/home/ubuntu/sqsblaster.py")

	cmd.Stdin = bytes.NewReader(payload)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Printf("sqsblaster.py failed %v", err)
	}
	return payload
}

// need three or two args . one is first run , bool . second is filetype
// eg curl -H "Authorization: Bearer mykey" http://localhost/seekAuto  -d 'first live'
func SeekExpiredAuto(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, r.Method, http.StatusBadRequest)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}

	body, err := io.ReadAll(r.Body)

	first_run_and_file_type := string(body)
	fmt.Println(first_run_and_file_type)
	first_run := false
	//var filetype *string
	//TODO dont use data for identification use something else
	var roles []string

	var filetype string
	if strings.Contains(first_run_and_file_type, "first") {
		first_run = true
		filetype = strings.Split(first_run_and_file_type, " ")[1]

		if err != nil {
			log.Println(err)
			http.Error(w, "Problem executing first run data could be unexpected format", http.StatusBadRequest)
		}
		roles, err = sqlconnect.SeekExpiredAuto(filetype, first_run)

	} else {
		filetype = first_run_and_file_type

		roles, err = sqlconnect.SeekExpiredAuto(filetype, first_run)
	}

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if roles == nil {
		fmt.Println("Job Done ")
		return
	}
	payload, _ := json.Marshal(roles)

	auto := "true"
	cmd := exec.Command("python3", "/home/ubuntu/backfill.py", filetype, auto)

	cmd.Stdin = bytes.NewReader(payload)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Printf("backfill.py failed %v", err)
	}

	response := struct {
		Status     string `json:"status"`
		StatusCode int
	}{
		Status:     fmt.Sprintf(" Successfully sent auto. have a good day %d", len(roles)),
		StatusCode: 200,
	}

	json.NewEncoder(w).Encode(response)
}

// should be expobackoff will pass fileid from lambda. for now only search term to update .
func Backoff(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, r.Method, http.StatusBadRequest)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "Problem reading could be unexpected format", http.StatusBadRequest)
		return
	}

	type BackoffSearchTerms struct {
		Search_term_id int `json:"search_term_id"`
		Timestamp      int `json:"timestamp"`
	}

	// if minus 1 then we know that it came from orchestrator . in which case reset the sql search term.
	// if timestamp though then use scheduler for next call.
	var Search_term BackoffSearchTerms

	err = json.Unmarshal(body, &Search_term)
	if err != nil {
		log.Println(err)
		http.Error(w, "Problem ummarshalling json could be unexpected format", http.StatusBadRequest)
		return
	}

	err = sqlconnect.BackoffUpdate(Search_term.Search_term_id)

	if err != nil {
		log.Println(err)
		http.Error(w, "Problem with db update on backoff searchtermid", http.StatusBadRequest)
		return
	}

	if Search_term.Timestamp != -1 {

		//do time.dofunc
		delay := time.Until(time.Unix(int64(Search_term.Timestamp), 0))

		time.AfterFunc(delay, func() {

			search_term, err := sqlconnect.GetSearchTerms(false, -1, "SEARCH_TERM")
			if err != nil {
				log.Println("error getting search terms", err)
				return
			}

			if search_term == nil {
				fmt.Println("Job Done either second last run or last run ")
				return
			}

			fmt.Println("here is your data", search_term, len(search_term))

			payload := SendToSqs(search_term)

			log.Printf("successfully blasted sqs, %d items sent", len(payload))
		})
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "scheduled",
	})

}

// for manual runs only todo make auto
func SeekScroller(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, r.Method, http.StatusBadRequest)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
	}
	defer r.Body.Close()

	search_term := string(body)

	cmd := exec.Command("python3", "/home/ubuntu/scroller.py", search_term)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("scroller.py failed %v", err)
	}

}

func SeekCompanyDeed(w http.ResponseWriter, r *http.Request) {

	FreshCompany, err := sqlconnect.SeekCompanyChecker()

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	payload, _ := json.Marshal(FreshCompany)

	cmd := exec.Command("python3", "/home/ubuntu/scroller.py")

	cmd.Stdin = bytes.NewReader(payload)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Printf("scroller.py failed %v", err)
	}

}

// TODO MAKE SCROOLERV2 a loop for sqlconnect.SearchTerm array
// func SendToScroller(SearchTerms []sqlconnect.SearchTerm) []byte {
func SendToScroller(SearchTerms []sqlconnect.SearchTerm) {
	//payload, _ := json.Marshal(SearchTerms)
	search_term := SearchTerms[0].Search_term
	search_term_id := strconv.Itoa(SearchTerms[0].Search_term_id)

	cmd := exec.Command("python3", "/home/ubuntu/scrollerv2.py", search_term, search_term_id)

	//cmd.Stdin = bytes.NewReader(payload)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Printf("scrollerv2.py failed %v", err)
	}
	//return payload
}

// func SendToScrollersqs(SearchTerms []sqlconnect.SearchTerm) {
// 	//payload, _ := json.Marshal(SearchTerms)

// 	search_term := SearchTerms[0].Search_term
// 	search_term_id := strconv.Itoa(SearchTerms[0].Search_term_id)

// 	cmd := exec.Command("python3", "/home/ubuntu/scrollerv3.py", search_term, search_term_id)

// 	//cmd.Stdin = bytes.NewReader(payload)
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr

// 	if err := cmd.Run(); err != nil {
// 		log.Printf("scrollerv3.py failed %v", err)
// 	}
// 	//return payload
// }

func SendToScrollersqs(SearchTerms []sqlconnect.SearchTerm) {
	//payload, _ := json.Marshal(SearchTerms)

	for _, term := range SearchTerms {
		search_term := term.Search_term
		search_term_id := strconv.Itoa(term.Search_term_id)
		cmd := exec.Command("python3", "/home/ubuntu/scrollerv3.py", search_term, search_term_id)

		//cmd.Stdin = bytes.NewReader(payload)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Printf("scrollerv3.py failed %v", err)
		}
		time.Sleep(3 * time.Second)
	}

	//return payload
}

// summon the elector counts! //change this to array and pass in as stdin
func SummonSpot(term sqlconnect.SearchTerm, instance_id string, firstrun string) error {

	search_term := term.Search_term
	search_term_id := strconv.Itoa(term.Search_term_id)

	if firstrun == "yes" {
		// then instance id is None
		cmd := exec.Command("python3", "/home/ubuntu/sendssm.py", "yes", search_term, search_term_id, "")

		if err := cmd.Run(); err != nil {
			return utils.ErrorHandler(err, "program did not execute")
		}

		return nil
	}
	cmd := exec.Command("python3", "/home/ubuntu/sendssm.py", "no", search_term, search_term_id, instance_id)

	if err := cmd.Run(); err != nil {
		return utils.ErrorHandler(err, "program did not execute")
	}

	return nil
}

func SummonSpotv2(SearchTerms []sqlconnect.SearchTerm, instance_id string, firstrun string) ([]byte, error) {

	payload, _ := json.Marshal(SearchTerms)

	if firstrun == "yes" {
		// then instance id is None
		cmd := exec.Command("python3", "/home/ubuntu/sendssm2.py", "yes", "")
		cmd.Stdin = bytes.NewReader(payload)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return nil, utils.ErrorHandler(err, "program did not execute")
		}

		return payload, nil
	}
	cmd := exec.Command("python3", "/home/ubuntu/sendssm2.py", "no", instance_id)
	cmd.Stdin = bytes.NewReader(payload)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, utils.ErrorHandler(err, "program did not execute")
	}

	return payload, nil
}

type BlastRequest struct {
	FirstRun       bool   `json:"first_run"`
	NumberAccounts int    `json:"number_accounts"`
	InstanceID     string `json:"instance_id,omitempty"`
}

// like deed blaster but for spot
func SpotDeedBlaster(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, r.Method, http.StatusBadRequest)
		http.Error(w, "method not allowed ", http.StatusMethodNotAllowed)
	}

	var req BlastRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println(err)
		http.Error(w, "Invalid Request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	SearchTerms, err := sqlconnect.GetSearchTerms(req.FirstRun, req.NumberAccounts, "SEARCH_TERM_DEED")

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if SearchTerms == nil {
		fmt.Println("Job Done either second last run or last run ")
		return
	}

	fmt.Println("here is your data", SearchTerms, len(SearchTerms))

	//SendToScroller(SearchTerms)
	//SendToScrollersqs(SearchTerms)
	firstRunStr := "no"
	if req.FirstRun {
		firstRunStr = "yes"
	}
	// if err := SummonSpot(SearchTerms[0], req.InstanceID, firstRunStr); err != nil {
	// 	log.Println(err)
	// 	http.Error(w, "Problem summoning spot instance", http.StatusInternalServerError)
	// 	return
	// }
	payload, err := SummonSpotv2(SearchTerms, req.InstanceID, firstRunStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Problem summoning spot instance", http.StatusInternalServerError)
		return
	}

	fmt.Println("blasted the spot", len(payload))

	w.Header().Set("Content-Type", "application/json")

	response := struct {
		Status     string `json:"status"`
		StatusCode int
	}{
		Status:     fmt.Sprintf(" Successfully blasted summonspot. have a good day"),
		StatusCode: 200,
	}

	json.NewEncoder(w).Encode(response)
}
func DeedBlaster(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, r.Method, http.StatusBadRequest)
		http.Error(w, "method not allowed ", http.StatusMethodNotAllowed)
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Invalid requestr body", http.StatusBadRequest)
	}
	defer r.Body.Close()

	first_run_and_number_accounts := string(body)
	first_run := false
	number_accounts := 0
	//TODO dont use data for identification use something else
	if strings.Contains(first_run_and_number_accounts, "first") {

		first_run = true
		number_accounts, err = strconv.Atoi(strings.Split(first_run_and_number_accounts, " ")[1])

		if err != nil {
			log.Println(err)
			http.Error(w, "Problem executing most likely data in unexpected format", http.StatusInternalServerError)
		}

	} else {
		number_accounts, err = strconv.Atoi(first_run_and_number_accounts)
		if err != nil {
			log.Println(err)
			http.Error(w, "Problem executing most likely data in unexpected format", http.StatusInternalServerError)
		}
	}

	fmt.Println(first_run_and_number_accounts)

	SearchTerms, err := sqlconnect.GetSearchTerms(first_run, number_accounts, "SEARCH_TERM_DEED")

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if SearchTerms == nil {
		fmt.Println("Job Done either second last run or last run ")
		return
	}

	fmt.Println("here is your data", SearchTerms, len(SearchTerms))

	//SendToScroller(SearchTerms)
	SendToScrollersqs(SearchTerms)

	w.Header().Set("Content-Type", "application/json")

	response := struct {
		Status     string `json:"status"`
		StatusCode int
	}{
		Status:     fmt.Sprintf(" Successfully blasted scroller. have a good day"),
		StatusCode: 200,
	}

	json.NewEncoder(w).Encode(response)

}

func SeekAutoCompany(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, r.Method, http.StatusBadRequest)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "Problem reading could be unexpected format", http.StatusBadRequest)
		return
	}
	// TODO make automatic
	if string(body) == "True" {
		//then its automatic

		response := struct {
			Status     string `json:"status"`
			StatusCode int
		}{
			Status:     fmt.Sprintf(" please implement"),
			StatusCode: 200,
		}

		json.NewEncoder(w).Encode(response)
	} else {

		var CompaniesDeets, err = sqlconnect.GetCompanyDeets()

		if err != nil {
			log.Println(err)
			http.Error(w, "Problem reading from db could be unexpected format", http.StatusInternalServerError)
			return
		}

		payload, _ := json.Marshal(CompaniesDeets)
		auto := "false"

		cmd := exec.Command("python3", "/home/ubuntu/details.py", auto)
		cmd.Stdin = bytes.NewReader(payload)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Printf("details.py failed %v", err)
		}
	}
}

func RedirectIndLinker(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, r.Method, http.StatusBadRequest)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "problem readinf could be unexpected format", http.StatusBadRequest)
	}

	if string(body) == "True" {
		//then its automatic

		response := struct {
			Status     string `json:"status"`
			StatusCode int
		}{
			Status:     fmt.Sprintf(" please implement"),
			StatusCode: 200,
		}

		json.NewEncoder(w).Encode(response)
	} else {

		var JobRedirect_DEED, err = sqlconnect.RedirectInd()

		if err != nil {
			log.Println(err)
			http.Error(w, "Problem reading from db could be unexpected format", http.StatusInternalServerError)
			return
		}

		payload, _ := json.Marshal(JobRedirect_DEED)
		auto := "false"

		cmd := exec.Command("python3", "/home/ubuntu/redirecter.py", auto)
		cmd.Stdin = bytes.NewReader(payload)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Printf("redirecter.py failed %v", err)
		}

	}

}
