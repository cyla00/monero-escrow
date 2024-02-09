package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/redis/go-redis/v9"
)

type Dbs struct {
	Psql  *sql.DB
	Redis *redis.Client
}

// handlers
func (db *Dbs) PostSignup(w http.ResponseWriter, r *http.Request) {

	type login struct {
		Username string
		Password string
	}
	var user_login login
	reqBody, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		http.Error(w, "error", http.StatusBadRequest)
		return
	}
	body := json.Unmarshal(reqBody, &user_login)
	json.NewEncoder(w).Encode(&user_login)

	fmt.Fprintf(w, "%+v", body)
}

func (db *Dbs) PostSignin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "login page")
}

func (db *Dbs) PostSellerContractOk(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "login page")
}

func (db *Dbs) PostBuyerInitTransaction(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "login page")
}

func (db *Dbs) PostBuyerTransactionOk(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "login page")
}
