package main

import (
	"log"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/cyla00/monero-escrow/middleware"
	"github.com/cyla00/monero-escrow/routes"
	"github.com/cyla00/monero-escrow/views"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	// ## static routes ##
	http.Handle("/", templ.Handler(views.Index()))
	http.Handle("/sign-up", templ.Handler(views.Signup()))
	http.Handle("/sign-in", templ.Handler(views.Signin()))
	http.Handle("/transaction", templ.Handler(views.Transaction())) // accepts query ?id=transaction-id

	// ## no AUTH routes ##

	// ## AUTH routes ##
	http.Handle("/api"+os.Getenv("API_VERSION")+"/sign-in", http.HandlerFunc(routes.PostSignin))
	http.Handle("/api"+os.Getenv("API_VERSION")+"/sign-up", http.HandlerFunc(routes.PostSignup))
	// AUTH buyer routes
	http.Handle("/api"+os.Getenv("API_VERSION")+"/buyer/init-transaction", middleware.AuthMiddleware(http.HandlerFunc(routes.PostBuyerInitTransaction)))       // create contract + deposit
	http.Handle("/api"+os.Getenv("API_VERSION")+"/buyer/transaction-confirmation", middleware.AuthMiddleware(http.HandlerFunc(routes.PostBuyerTransactionOk))) // buyer confirms
	// AUTH seller routes
	http.Handle("/api"+os.Getenv("API_VERSION")+"/seller/verify-contract", middleware.AuthMiddleware(http.HandlerFunc(routes.PostSellerContractOk))) // verify contract (yes/no) + 10% hostage deposit

	log.Print("http://127.0.0.1:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
