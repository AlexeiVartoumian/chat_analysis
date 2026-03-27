package router

import (
	"api/handlers"
	"net/http"
)

func MainRouter() *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /lastThreeDays", handlers.GetLastThreeDays)

	return mux

}
