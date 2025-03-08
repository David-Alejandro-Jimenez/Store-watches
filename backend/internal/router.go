package internal

import (
	"net/http"

	"github.com/David-Alejandro-Jimenez/sale-watches/internal/handlers/private"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/handlers/public"
	ratelimiter "github.com/David-Alejandro-Jimenez/sale-watches/pkg/security/rate_limiter"
	securityAuth "github.com/David-Alejandro-Jimenez/sale-watches/pkg/security/security_auth"
	"github.com/gorilla/mux"
)

// The SetupRouter function is responsible for configuring and returning the main router of the web application, using the mux package to define the routes and corresponding handlers.
// 1. Static route configuration: Allows serving resource files (CSS, JS, images) directly from the frontend directory.
// 2. Definition of public routes: Routes accessible to any user, with rate limiting protection to prevent abuse.
// 3. Definition of protected routes: Routes that require authentication, combining rate limiting and authentication middlewares.
// 4. Main router: A mux.Router object is returned that centralizes all configured routes.
// This framework facilitates centralized route management and enforcement of security and performance policies (such as authentication and rate control) in a modular and scalable manner.
func SetupRouter() *mux.Router {
	var router = mux.NewRouter()

	//Use de javascript, html, css and images
	var staticDir = "./../frontend"
	router.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir(staticDir+"/css/"))))
	router.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir(staticDir+"/js/"))))
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir(staticDir+"/assets/"))))

	//Routes public
	router.Handle("/", ratelimiter.RateLimitMiddleware(http.HandlerFunc(public.Main_page))).Methods("GET")
	router.Handle("/register", ratelimiter.RateLimitMiddleware(http.HandlerFunc(public.RegisterPOST))).Methods("POST")
	router.Handle("/login", ratelimiter.RateLimitMiddleware(http.HandlerFunc(public.LoginPOST))).Methods("POST")
	router.Handle("/comments",  ratelimiter.RateLimitMiddleware(http.HandlerFunc(public.Comments))).Methods("GET")

	//Routes protected
	protectedHandler := securityAuth.AuthMiddleware(ratelimiter.RateLimitMiddleware(http.HandlerFunc(private.NewComment)))
	router.Handle("/comments/NewComment", protectedHandler).Methods("POST")

	return router
}