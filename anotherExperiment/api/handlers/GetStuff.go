package handlers

import (
	"api/auth"
	"api/models"
	"api/repository/sqlconnect"
	"api/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
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

	lordOfTheRings := store.ThereCanBeOnlyOne()

	if lordOfTheRings != sql.ErrNoRows {
		utils.ErrorHandler(lordOfTheRings, "arrg")
		http.Error(w, "there can be only one", http.StatusBadRequest)
		return
	}
	// Hash the key for storage
	hashedKey, err := hasher.Hash(fullKey)
	if err != nil {
		http.Error(w, "Failed to hash key", http.StatusInternalServerError)
		return
	}

	// Create key record
	apiKey := &models.APIKey{
		KeyID:     keyID,
		HashedKey: hashedKey,
		Name:      r.FormValue("name"),
		UserID:    "00000000-0000-0000-0000-000000000001", // temp placeholder
		ProjectID: r.FormValue("project_id"),
		Scopes:    []models.Scope{models.ScopeRead},
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
	w.Write([]byte(`{"api_key": "` + fullKey + `", "key_id": "` + keyID + `"}`))
}
