package internal

import (
	"net/http"

	"github.com/David-Alejandro-Jimenez/venta-relojes/internal/handlers/private"
	"github.com/David-Alejandro-Jimenez/venta-relojes/internal/handlers/public"
	"github.com/David-Alejandro-Jimenez/venta-relojes/pkg/security"
	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	var router = mux.NewRouter()

	//Use de javascript, html, css and images
	var staticDir = "./../frontend"
	router.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir(staticDir+"/css/"))))
	router.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir(staticDir+"/js/"))))
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir(staticDir+"/assets/"))))

	//Routes public
	router.HandleFunc("/", public.Main_page).Methods("GET")
	router.HandleFunc("/register", public.RegisterPOST).Methods("POST")
	router.HandleFunc("/login", public.LoginPOST).Methods("POST")
	router.HandleFunc("/comments", public.Comments).Methods("GET")

	//Routes protected
	router.Handle("/comments/NewComment", security.AuthMiddleware(http.HandlerFunc(private.NewComment))).Methods("POST")

	return router
}