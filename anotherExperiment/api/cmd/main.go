package main

import (
	"api/auth"
	"api/repository/sqlconnect"
	"api/router"
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

func main() {

	fmt.Println("testing sql connection")

	db, err := sqlconnect.ConnectDb()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	port := "443"

	store := sqlconnect.NewPostgresStore(db)

	if err := store.CreateTable(context.Background()); err != nil {
		log.Fatal("Failed to create tables:", err)
	}
	generator := auth.NewAPIKeyGenerator()
	hasher := auth.NewKeyHasher()
	authMiddleware := auth.NewAuthMiddleware(generator, hasher, store)
	router := router.MainRouter(authMiddleware)

	//timeout againt malicious clients
	server := &http.Server{
		Addr: "0.0.0.0:" + port,
		//Addr:    "0.0.0.0:" + port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	hostPolicy := func(ctx context.Context, host string) error {
		allowedHost := "jobdice.worldcaffeine.com"

		if host == allowedHost {
			return nil
		}
		return fmt.Errorf("acme/autocert: only %s host is allowed", allowedHost)
	}

	dataDir := "."

	m := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: hostPolicy,
		Cache:      autocert.DirCache(dataDir),
	}

	//err = server.ListenAndServe()

	server.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}

	go func() {
		err := server.ListenAndServeTLS("", "")
		if err != nil {
			log.Fatalln(err, "error starting the server")
		}
	}()
	fmt.Println("server is up and running on port", port)

	//need a http server at port 80 for letsencrypt + autocert
	httpSrv := &http.Server{
		Addr: ":80",
		Handler: m.HTTPHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "https://"+r.Host+r.URL.String(), http.StatusMovedPermanently)
		})),
	}
	err = httpSrv.ListenAndServe()
	if err != nil {
		log.Fatalf("httpSrv.ListenAndServe() failed with %s", err)
	}

	//sqlconnect.BackfillEmbeddings()

	//api stuff
	//sqlconnect.SearchSimilarJobs("react developer")

}
