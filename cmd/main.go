package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/cyla00/monero-escrow/routes"
	"github.com/cyla00/monero-escrow/views"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

var PSQL *sql.DB
var REDIS *redis.Client

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

	dbInject := routes.Dbs{
		Psql:  PSQL,
		Redis: REDIS,
	}

	// ## static routes ##
	http.Handle("/", templ.Handler(views.Index()))
	http.Handle("/sign-up", templ.Handler(views.Signup()))
	http.Handle("/sign-in", templ.Handler(views.Signin()))
	http.Handle("/transaction", templ.Handler(views.Transaction())) // accepts query ?id=transaction-id

	// ## no AUTH routes ##

	// ## AUTH routes ##
	http.Handle("/api"+os.Getenv("API_VERSION")+"/sign-in", http.HandlerFunc(dbInject.PostSignin))
	http.Handle("/api"+os.Getenv("API_VERSION")+"/sign-up", http.HandlerFunc(dbInject.PostSignup))

	// AUTH buyer routes
	http.Handle("/api"+os.Getenv("API_VERSION")+"/buyer/init-transaction", authMiddleware(http.HandlerFunc(dbInject.PostBuyerInitTransaction)))       // create contract + deposit
	http.Handle("/api"+os.Getenv("API_VERSION")+"/buyer/transaction-confirmation", authMiddleware(http.HandlerFunc(dbInject.PostBuyerTransactionOk))) // buyer confirms

	// AUTH seller routes
	http.Handle("/api"+os.Getenv("API_VERSION")+"/seller/verify-contract", authMiddleware(http.HandlerFunc(dbInject.PostSellerContractOk))) // verify contract (yes/no) + 10% hostage deposit

	log.Print("http://127.0.0.1:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

// authentication middleware
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print("auth middleware executed")
		next.ServeHTTP(w, r)
	})
}
