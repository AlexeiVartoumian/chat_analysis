package main

import (
	"api/router"
	"fmt"
	"log"
	"net/http"
)

func main() {

	fmt.Println("testing sql connection")

	//sqlconnect.CsvFile("C:/Users/wwwal/Downloads/processedJobs.csv", "JOBS")

	//fmt.Println(sqlconnect.GetJobById(1))

	port := "3000"

	router := router.MainRouter()

	server := &http.Server{
		Addr:    "localhost:" + port,
		Handler: router,
	}

	fmt.Println("server is up and running on port", port)

	err := server.ListenAndServe()

	if err != nil {
		log.Fatalln(err, "error starting the server")
	}
	// sqlconnect.CsvFile("C:/Users/wwwal/Downloads/processedJobs.csv", "COMPANY")
	// sqlconnect.CsvFile("C:/Users/wwwal/Downloads/processedJobs.csv", "JOBS")
	// sqlconnect.CsvFile("C:/Users/wwwal/Downloads/company_data.csv", "COMPANY_METADATA")

	// sqlconnect.CsvFile("C:/Users/wwwal/Downloads/job_metadata.csv", "JOB_METADATA")
	// sqlconnect.CsvFile("C:/Users/wwwal/Downloads/job_description.csv", "JOB_DESCRIPTION")

	//sqlconnect.BackfillEmbeddings()

	//api stuff
	//sqlconnect.SearchSimilarJobs("react developer")

}
