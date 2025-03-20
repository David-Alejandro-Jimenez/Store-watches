package internal

import (
	"database/sql"
	"net/http"

	//"github.com/David-Alejandro-Jimenez/sale-watches/internal/handlers/private"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/handlers/public"
	commentsRepository "github.com/David-Alejandro-Jimenez/sale-watches/internal/repository/comments"
	"github.com/David-Alejandro-Jimenez/sale-watches/pkg/security/rate_limiter"

	//securityAuth "github.com/David-Alejandro-Jimenez/sale-watches/pkg/security/security_auth"
	"github.com/gorilla/mux"
)

type RouterConfiguration interface {
	SetupRoutes(router *mux.Router)
}

type defaultRouterConfiguration struct {
	IPExtractor ratelimiter.IPExtractor
	RateLimiter ratelimiter.RateLimiterHandler
	HandlerComment public.HandlerComment
}

func (c *defaultRouterConfiguration) SetupRoutes(router *mux.Router) {
	var staticDir = "./../frontend"
	router.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir(staticDir+"/css/"))))
	router.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir(staticDir+"/js/"))))
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir(staticDir+"/assets/"))))

	router.Handle("/", ratelimiter.RateLimitMiddleware(http.HandlerFunc(public.MainPage), c.IPExtractor, c.RateLimiter)).Methods("GET")
	router.Handle("/register", ratelimiter.RateLimitMiddleware(http.HandlerFunc(public.RegisterPOST), c.IPExtractor, c.RateLimiter)).Methods("POST")
	router.Handle("/login", ratelimiter.RateLimitMiddleware(http.HandlerFunc(public.LoginPOST), c.IPExtractor, c.RateLimiter)).Methods("POST")
	router.Handle("/comments", ratelimiter.RateLimitMiddleware(http.HandlerFunc(c.HandlerComment.Comments), c.IPExtractor, c.RateLimiter)).Methods("GET")
}

func SetupRouter(db *sql.DB, rateHandler ratelimiter.RateLimiterHandler) *mux.Router {
	router := mux.NewRouter()
	commentsRepo := commentsRepository.NewComments(db)
	
	publicHandlerComment := public.NewHandlerComment(commentsRepo)
	configurator := &defaultRouterConfiguration{
		IPExtractor: &ratelimiter.DefaultIPExtractor{},
		RateLimiter: rateHandler,
		HandlerComment: *publicHandlerComment,
	}
	configurator.SetupRoutes(router)
	return router
}