package router

import (
	"api/auth"
	"api/handlers"
	"api/models"
	"net/http"
)

// func MainRouter() *http.ServeMux {

// 	mux := http.NewServeMux()

// 	mux.HandleFunc("GET /lastThreeDays", handlers.GetLastThreeDays)

// 	authRouter := AuthRouter()

// 	mux.Handle("/", authRouter)
// 	return mux

// }

func MainRouter(authMiddleware *auth.AuthMiddleware) *http.ServeMux {
	mux := http.NewServeMux()

	// Public — generate a key
	mux.HandleFunc("POST /api/keys", handlers.PostApiKey)

	// Protected — requires valid API key
	mux.Handle("GET /lastThreeDays", authMiddleware.Authenticate(models.ScopeRead)(
		http.HandlerFunc(handlers.GetLastThreeDays),
	))

	mux.Handle("POST /semanticSearch", authMiddleware.Authenticate(models.ScopeRead)(
		http.HandlerFunc(handlers.SemanticSearch),
	))

	mux.Handle("GET /onlyCompanyLinks", authMiddleware.Authenticate(models.ScopeRead)(
		http.HandlerFunc(handlers.CompanyUrlOnly),
	))
	return mux
}
