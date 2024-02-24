package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/cyla00/monero-escrow/routes"
	"github.com/cyla00/monero-escrow/views"
	"github.com/icholy/digest"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	PSQL, sql_err := sql.Open("postgres", os.Getenv("PSQL_CONN_URL"))

	if sql_err != nil {
		log.Fatal(sql_err)
	}

	REDIS := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: os.Getenv("REDIS_PWD"),
		DB:       0,
	})

	xmrAuthClient := &http.Client{
		Transport: &digest.Transport{
			Username: os.Getenv("XMR_USER"),
			Password: os.Getenv("XMR_PWD"),
		},
	}

	Inject := routes.Injection{
		Psql:          PSQL,
		Redis:         REDIS,
		XmrAuthClient: xmrAuthClient,
	}

	var baseApiUrl = fmt.Sprintf("/api/%s", os.Getenv("API_VERSION"))

	// ## static routes ##
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/", routes.GetRequestMiddleware(http.HandlerFunc(Inject.GetIndexView)))
	http.Handle("/404", routes.GetRequestMiddleware(templ.Handler(views.NotFound(), templ.WithStatus(http.StatusNotFound))))
	http.Handle("/sign-up", routes.GetRequestMiddleware(http.HandlerFunc(Inject.GetSignupView)))
	http.Handle("/sign-in", routes.GetRequestMiddleware(http.HandlerFunc(Inject.GetSigninView)))
	http.Handle("/transaction", routes.GetRequestMiddleware(http.HandlerFunc(Inject.GetTransactionPayment))) // accepts query ?id=transaction-id

	// ## no AUTH routes ##

	// ## API routes ##
	http.Handle(baseApiUrl+"/sign-in", routes.PostRequestMiddleware(http.HandlerFunc(Inject.PostSignin)))
	http.Handle(baseApiUrl+"/sign-up", routes.PostRequestMiddleware(http.HandlerFunc(Inject.PostSignup)))
	http.Handle(baseApiUrl+"/reset-password", routes.PutRequestMiddleware(http.HandlerFunc(Inject.PostChangePassword)))

	// AUTH buyer routes
	http.Handle(baseApiUrl+"/buyer/init-transaction", routes.PostRequestMiddleware(Inject.AuthMiddleware(http.HandlerFunc(Inject.PostBuyerInitTransaction))))       // create contract + deposit
	http.Handle(baseApiUrl+"/buyer/transaction-confirmation", routes.PostRequestMiddleware(Inject.AuthMiddleware(http.HandlerFunc(Inject.PostBuyerTransactionOk)))) // buyer confirms

	// AUTH seller routes
	http.Handle(baseApiUrl+"/seller/verify-contract", routes.PostRequestMiddleware(Inject.AuthMiddleware(http.HandlerFunc(Inject.PostSellerContractOk)))) // verify contract (yes/no) + 10% hostage deposit

	log.Print("http://127.0.0.1:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
