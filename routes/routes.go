package routes

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/cyla00/monero-escrow/passwords"
	"github.com/cyla00/monero-escrow/types"
	"github.com/fossoreslp/go-uuid-v4"
	"github.com/redis/go-redis/v9"
)

type Dbs struct {
	Psql  *sql.DB
	Redis *redis.Client
}

// HANDLERS
func (db *Dbs) PostSignup(w http.ResponseWriter, r *http.Request) {

	type login struct {
		Username string
		Password string
	}
	var body *login
	reqBody, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errMsg := types.JsonResponse{
			Message: "Error, please retry later",
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}
	json.Unmarshal(reqBody, &body)

	if body.Username == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errMsg := types.JsonResponse{
			Message: "Provide a username",
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	if body.Password == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errMsg := types.JsonResponse{
			Message: "Provide a password",
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	size, number, upper, special := passwords.VerifyPassword(body.Password)
	if !size || !number || !upper || !special {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errMsg := types.JsonResponse{
			Message: "Password too weak",
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}
	hashedPassword, newSalt, hashingError := passwords.HashPassword(body.Password)
	if hashingError != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errMsg := types.JsonResponse{
			Message: "An error occurred, retry later",
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}
	var queryResult types.User
	userCheckErr := db.Psql.QueryRow(
		"INSERT INTO users (username, password, salt) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING RETURNING *;",
		&body.Username, hashedPassword, newSalt).Scan(&queryResult.Id, &queryResult.Hash, &queryResult.Username, &queryResult.Password, &queryResult.Salt)
	if userCheckErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errMsg := types.JsonResponse{
			Message: "Username not available",
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	succMsg := types.JsonResponse{
		Message: "Successfully registered",
	}
	json.NewEncoder(w).Encode(succMsg)
}

func (db *Dbs) PostSignin(w http.ResponseWriter, r *http.Request) {
	username, password, _ := r.BasicAuth()
	if username == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errMsg := types.JsonResponse{
			Message: "Provide your username",
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}
	if password == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errMsg := types.JsonResponse{
			Message: "Provide your password",
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	var queryResult types.User
	queryErr := db.Psql.QueryRow(
		"SELECT * FROM users WHERE username=$1", username).Scan(&queryResult.Id, &queryResult.Hash, &queryResult.Username, &queryResult.Password, &queryResult.Salt)
	if queryErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errMsg := types.JsonResponse{
			Message: "Incorrect credentials",
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}
	passCheck := passwords.CheckPasswords(password, queryResult.Password, queryResult.Salt)
	if !passCheck {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		errMsg := types.JsonResponse{
			Message: "Incorrect credentials",
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	ctx := context.Background()
	newSessionId, _ := uuid.NewString()
	hashedSessionId := passwords.Hash256(newSessionId)
	sessionError := db.Redis.Set(ctx, hashedSessionId, queryResult.Id, 168*time.Hour).Err()
	if sessionError != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errMsg := types.JsonResponse{
			Message: "Error, retry later",
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	authCookie := http.Cookie{
		Name:     "fidexmr",
		Value:    hashedSessionId,
		MaxAge:   3600 * 5,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &authCookie)
	w.WriteHeader(http.StatusOK)
	succMsg := types.JsonResponse{
		Message: "Connected",
	}
	json.NewEncoder(w).Encode(succMsg)
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

// MIDDLEWARES
func (db *Dbs) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, cookieErr := r.Cookie("fidexmr")
		if cookieErr != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.Background()
		_, getErr := db.Redis.Get(ctx, cookie.Value).Result()
		if getErr != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		log.Print("auth middleware executed")
		next.ServeHTTP(w, r)
	})
}

func PostRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "bad request", http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func PutRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			http.Error(w, "bad request", http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func DeleteRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			http.Error(w, "bad request", http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	})
}
