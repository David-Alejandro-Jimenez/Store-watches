package internal

import (
	"net/http"

	//"github.com/David-Alejandro-Jimenez/sale-watches/internal/handlers/private"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/handlers/public"
	"github.com/David-Alejandro-Jimenez/sale-watches/pkg/security/rate_limiter"
	//securityAuth "github.com/David-Alejandro-Jimenez/sale-watches/pkg/security/security_auth"
	"github.com/gorilla/mux"
)

type RouterConfiguration interface {
	SetupRoutes(router *mux.Router)
}

type DefaultRouterConfiguration struct {
	IPExtractor ratelimiter.IPExtractor
	RateLimiter ratelimiter.RateLimiterHandler
}

func (c *DefaultRouterConfiguration) SetupRoutes(router *mux.Router) {
	var staticDir = "./../frontend"
	router.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir(staticDir+"/css/"))))
	router.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir(staticDir+"/js/"))))
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir(staticDir+"/assets/"))))

	router.Handle("/", ratelimiter.RateLimitMiddleware(http.HandlerFunc(public.Main_page), c.IPExtractor, c.RateLimiter)).Methods("GET")
	router.Handle("/register", ratelimiter.RateLimitMiddleware(http.HandlerFunc(public.RegisterPOST), c.IPExtractor, c.RateLimiter)).Methods("POST")
	router.Handle("/login", ratelimiter.RateLimitMiddleware(http.HandlerFunc(public.LoginPOST), c.IPExtractor, c.RateLimiter)).Methods("POST")
	router.Handle("/comments", ratelimiter.RateLimitMiddleware(http.HandlerFunc(public.Comments), c.IPExtractor, c.RateLimiter)).Methods("GET")
}

func SetupRouter(rateHandler ratelimiter.RateLimiterHandler) *mux.Router {
	router := mux.NewRouter()
	configurator := &DefaultRouterConfiguration{
		IPExtractor: &ratelimiter.DefaultIPExtractor{},
		RateLimiter: rateHandler,
	}
	configurator.SetupRoutes(router)
	return router
}