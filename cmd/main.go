package main

import (
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/cyla00/monero-escrow/middleware"
	"github.com/cyla00/monero-escrow/routes"
	"github.com/cyla00/monero-escrow/views"
)

func main() {
	// statuc GET routes
	http.Handle("/", templ.Handler(views.Index()))
	http.Handle("/sign-up", templ.Handler(views.Signup()))
	http.Handle("/sign-in", templ.Handler(views.Signin()))

	// no AUTH routes

	// AUTH routes

	// authRoutes := http.NewServeMux()
	http.Handle("/test", middleware.AuthMiddleware(http.HandlerFunc(routes.GetTest)))

	log.Print("http://127.0.0.1:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
