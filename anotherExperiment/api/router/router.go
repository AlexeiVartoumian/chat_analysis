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

	mux.Handle("POST /semanticSearch", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(handlers.SemanticSearch),
	))

	mux.Handle("GET /sqsBlaster", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(handlers.SqsBlaster),
	))

	mux.Handle("GET /onlyCompanyLinks", authMiddleware.Authenticate(models.ScopeRead)(
		http.HandlerFunc(handlers.CompanyUrlOnly),
	))

	mux.Handle("GET /seekExpired", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(handlers.SeekExpiredRoles),
	))

	mux.Handle("GET /seekReopened", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(handlers.SeekReopenedRoles),
	))

	mux.Handle("GET /insertDb", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(handlers.InsertToDb),
	))

	mux.Handle("POST /queryDb", authMiddleware.Authenticate(models.ScopeRead)(
		http.HandlerFunc(handlers.HandleQuery),
	))

	return mux
}
