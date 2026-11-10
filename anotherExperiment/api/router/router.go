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

	mux.Handle("POST /sendlevCompany", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SendlevCompany),
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

	//athena
	mux.Handle("POST /seekGreen", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SeekGreen),
	))
	//zoos zoos
	mux.Handle("POST /seekAsh", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SeekAsh),
	))

	mux.Handle("POST /seekAshJd", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SeekAshJd),
	))

	mux.Handle("POST /seekdeedJd", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SeekDeedJd),
	))

	mux.Handle("POST /sendworkweek", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SendWorkweek),
	))

	mux.Handle("POST /sendlev", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SendLev),
	))

	mux.Handle("POST /sendworkweekJd", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SendWorkweekJd),
	))

	mux.Handle("POST /sendgreenlifecycle", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SendGreenLifecycle),
	))

	mux.Handle("POST /sendworklifecycle", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SendWorkLifecycle),
	))

	mux.Handle("POST /sendorac", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SendOrac),
	))

	mux.Handle("POST /sendoraclifecycle", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SendOracLifecycle),
	))

	mux.Handle("POST /sendoraclifecyclegcp", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SendOracLifecycle),
	))

	mux.Handle("POST /sendoracJd", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SendOracJd),
	))

	mux.Handle("POST /sendsmart", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SendSmart),
	))

	mux.Handle("POST /sendsmartJd", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SendSmartJD),
	))

	mux.Handle("POST /sendsmartlifecycle", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SendSmartLifecycle),
	))

	mux.Handle("POST /sendbambi", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SendBambi),
	))

	mux.Handle("POST /sendbambiJd", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SendBambiJD),
	))

	mux.Handle("POST /sendbambilifecycle", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SendBambiLifecycle),
	))

	mux.Handle("POST /sendworkaworka", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SendWorkaworka),
	))

	mux.Handle("POST /sendworkaworkaJd", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SendWorkaworkaJD),
	))

	mux.Handle("POST /sendworkaworkalifecycle", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SendWorkaworkaLifecycle),
	))

	mux.Handle("POST /seekRip", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SendRip),
	))
	mux.Handle("POST /sendriplifecycle", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SendRipLifecycle),
	))

	mux.Handle("POST /seekRecruit", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SendRecruit),
	))
	mux.Handle("POST /sendrecruitlifecycle", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SendRecruitLifecycle),
	))

	mux.Handle("POST /seekBrez", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SendBrez),
	))
	mux.Handle("POST /sendbrezlifecycle", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SendBrezLifecycle),
	))

	mux.Handle("POST /seekTailor", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SendTailor),
	))
	mux.Handle("POST /sendtailorlifecycle", authMiddleware.Authenticate(models.ScopeAdmin)(
		http.HandlerFunc(h.SendTailorLifecycle),
	))

	return mux
}
