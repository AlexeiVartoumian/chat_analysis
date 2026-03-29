package main

import (
	"api/auth"
	"api/repository/sqlconnect"
	"api/router"
	"context"
	"fmt"
	"log"
	"net/http"
)

func main() {

	//sqlconnect.CsvFile("C:/Users/wwwal/Downloads/processedJobs.csv", "JOBS")

	//fmt.Println(sqlconnect.GetJobById(1))
	fmt.Println("testing sql connection")

	db, err := sqlconnect.ConnectDb()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	port := "3000"

	store := sqlconnect.NewPostgresStore(db)

	if err := store.CreateTable(context.Background()); err != nil {
		log.Fatal("Failed to create tables:", err)
	}
	generator := auth.NewAPIKeyGenerator()
	hasher := auth.NewKeyHasher()
	authMiddleware := auth.NewAuthMiddleware(generator, hasher, store)
	router := router.MainRouter(authMiddleware)

	server := &http.Server{
		Addr: "localhost:" + port,
		//Addr:    "0.0.0.0:" + port,
		Handler: router,
	}

	fmt.Println("server is up and running on port", port)

	err = server.ListenAndServe()

	if err != nil {
		log.Fatalln(err, "error starting the server")
	}
	// sqlconnect.CsvFile("C:/Users/wwwal/Downloads/processedJobs.csv", "COMPANY")
	// sqlconnect.CsvFile("C:/Users/wwwal/Downloads/processedJobs.csv", "JOBS")
	// sqlconnect.CsvFile("C:/Users/wwwal/Downloads/company_data.csv", "COMPANY_METADATA")

	// sqlconnect.CsvFile("C:/Users/wwwal/Downloads/job_metadata.csv", "JOB_METADATA")
	// sqlconnect.CsvFile("C:/Users/wwwal/Downloads/job_description.csv", "JOB_DESCRIPTION")

	// sqlconnect.CsvFile("/home/ubuntu/processedJobs.csv", "COMPANY")
	//     sqlconnect.CsvFile("/home/ubuntu/processedJobs.csv", "JOBS")
	//     sqlconnect.CsvFile("/home/ubuntu/company_data.csv", "COMPANY_METADATA")

	//     sqlconnect.CsvFile("/home/ubuntu/job_metadata.csv", "JOB_METADATA")
	//     sqlconnect.CsvFile("/home/ubuntu/job_description.csv", "JOB_DESCRIPTION")

	//sqlconnect.BackfillEmbeddings()

	//api stuff
	//sqlconnect.SearchSimilarJobs("react developer")

}

// package main

// import (
//         "api/router"
//         "fmt"
//         "log"
//         "net/http"
// //"api/repository/sqlconnect"
// )

// func main() {

//         fmt.Println("testing sql connection")

//         //sqlconnect.CsvFile("C:/Users/wwwal/Downloads/processedJobs.csv", "JOBS")

//         //fmt.Println(sqlconnect.GetJobById(1))

//         port := "3000"

//         router := router.MainRouter()

//         server := &http.Server{
//                 //Addr:    "0.0.0.0:" + port,
// 				Addr: "localhost:" + port,
//                 Handler: router,
//         }

//         fmt.Println("server is up and running on port", port)

//         err := server.ListenAndServe()

//         if err != nil {
//                 log.Fatalln(err, "error starting the server")
//         }
// //         //sqlconnect.CsvFile("/home/ubuntu/processedJobs.csv", "COMPANY")
// //         //sqlconnect.CsvFile("/home/ubuntu/processedJobs.csv", "JOBS")
// //         //sqlconnect.CsvFile("/home/ubuntu/company_data.csv", "COMPANY_METADATA")

// //         //sqlconnect.CsvFile("/home/ubuntu/job_metadata.csv", "JOB_METADATA")
// //         //sqlconnect.CsvFile("/home/ubuntu/job_description.csv", "JOB_DESCRIPTION")

// //         //sqlconnect.BackfillEmbeddings()

// //         //api stuff
// //         //sqlconnect.SearchSimilarJobs("react developer")

// }
