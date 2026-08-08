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

func MainRouter(authMiddleware *auth.AuthMiddleware, h *handlers.Handler) *http.ServeMux {
	mux := http.NewServeMux()

	// Public — generate a key
	mux.HandleFunc("POST /api/keys", h.PostApiKey)

	// Protected — requires valid API key
	mux.Handle("GET /lastThreeDays", authMiddleware.Authenticate(models.ScopeRead)(
		http.HandlerFunc(handlers.GetLastThreeDays),
	))

	mux.Handle("POST /semanticSearch", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SemanticSearch),
	))

	mux.Handle("POST /sqsBlaster", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SqsBlaster),
	))

	mux.Handle("POST /Backoff", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SqsBlaster),
	))

	mux.Handle("GET /onlyCompanyLinks", authMiddleware.Authenticate(models.ScopeRead)(
		http.HandlerFunc(h.CompanyUrlOnly),
	))

	mux.Handle("GET /seekExpired", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SeekExpiredRoles),
	))

	mux.Handle("POST /seekAuto", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SeekExpiredAuto),
	))

	mux.Handle("GET /seekReopened", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SeekReopenedRoles),
	))

	mux.Handle("GET /insertDb", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.InsertToDb),
	))

	mux.Handle("POST /queryDb", authMiddleware.Authenticate(models.ScopeRead)(
		http.HandlerFunc(h.HandleQuery),
	))

	//do i even need ?
	mux.Handle("POST /scroller", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(handlers.SeekScroller),
	))

	mux.Handle("POST /seekcompanydeed", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SeekCompanyDeed),
	))

	mux.Handle("POST /deedblaster", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.DeedBlaster),
	))

	mux.Handle("POST /spotdeed", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SpotDeedBlaster),
	))

	mux.Handle("POST /seekAutoCompany", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SeekAutoCompany),
	))

	mux.Handle("POST /redirectInd", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.RedirectIndLinkerAuto),
	))

	mux.Handle("POST /seekExpiredAutoDeed", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SeekExpiredAutoDeed),
	))

	mux.Handle("POST /seekAshLead", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SeekAshLead),
	))

	mux.Handle("POST /seekAshCompany", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SeekAshCompany),
	))

	mux.Handle("POST /seekGreenCompany", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SeekGreenCompany),
	))

	mux.Handle("POST /seekGreenLead", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SeekGreenLead),
	))

	mux.Handle("POST /SeekAshLeadLink", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SeekAshLeadLink),
	))

	mux.Handle("POST /seekGreenLeadLink", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SeekGreenLeadLink),
	))

	mux.Handle("POST /seekGreenJd", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SeekGreenJd),
	))

	mux.Handle("POST /seekAshJd", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SeekAshJd),
	))

	return mux
}
