package routes

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
)

type Dbs struct {
	Psql  *sql.DB
	Redis *redis.Client
}

func checkPostMethod(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(res, "bad request", http.StatusBadRequest)
		return
	}
}

func (db *Dbs) PostSignup(w http.ResponseWriter, r *http.Request) {
	checkPostMethod(w, r)
	fmt.Fprintf(w, "login page")
}

func (db *Dbs) PostSignin(w http.ResponseWriter, r *http.Request) {
	checkPostMethod(w, r)
	fmt.Fprintf(w, "login page")
}

func (db *Dbs) PostSellerContractOk(w http.ResponseWriter, r *http.Request) {
	checkPostMethod(w, r)
	fmt.Fprintf(w, "login page")
}

func (db *Dbs) PostBuyerInitTransaction(w http.ResponseWriter, r *http.Request) {
	checkPostMethod(w, r)
	fmt.Fprintf(w, "login page")
}

func (db *Dbs) PostBuyerTransactionOk(w http.ResponseWriter, r *http.Request) {
	checkPostMethod(w, r)
	fmt.Fprintf(w, "login page")
}
