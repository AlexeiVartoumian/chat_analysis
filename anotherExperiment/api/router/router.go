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

	mux.Handle("POST /sqsBlaster", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(handlers.SqsBlaster),
	))

	mux.Handle("POST /Backoff", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(handlers.SqsBlaster),
	))

	mux.Handle("GET /onlyCompanyLinks", authMiddleware.Authenticate(models.ScopeRead)(
		http.HandlerFunc(handlers.CompanyUrlOnly),
	))

	mux.Handle("GET /seekExpired", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(handlers.SeekExpiredRoles),
	))

	mux.Handle("POST /seekAuto", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(handlers.SeekExpiredAuto),
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

	mux.Handle("POST /scroller", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(handlers.SeekScroller),
	))

	mux.Handle("POST /seekcompanydeed", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(handlers.SeekCompanyDeed),
	))

	mux.Handle("POST /deedblaster", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(handlers.DeedBlaster),
	))

	mux.Handle("POST /spotdeed", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(handlers.SpotDeedBlaster),
	))

	mux.Handle("POST /seekAutoCompany", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(handlers.SeekAutoCompany),
	))

	mux.Handle("POST /redirectInd", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(handlers.RedirectIndLinkerAuto),
	))

	mux.Handle("POST /seekExpiredAutoDeed", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(handlers.SeekExpiredAutoDeed),
	))

	mux.Handle("POST /seekAshLead", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(handlers.SeekAshLead),
	))

	mux.Handle("POST /seekAshCompany", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(handlers.SeekAshCompany),
	))

	mux.Handle("POST /seekGreenLead", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(handlers.SeekGreenLead),
	))

	mux.Handle("POST /SeekAshLeadLink", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(handlers.SeekAshLeadLink),
	))

	mux.Handle("POST /seekGreenLeadLink", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(handlers.SeekGreenLeadLink),
	))

	return mux
}
