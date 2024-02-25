package routes

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	moneroapi "github.com/cyla00/monero-escrow/monero-api"
	"github.com/cyla00/monero-escrow/passwords"
	"github.com/cyla00/monero-escrow/types"
	"github.com/cyla00/monero-escrow/views"
	"github.com/fossoreslp/go-uuid-v4"
	"github.com/redis/go-redis/v9"
)

type Injection struct {
	Psql          *sql.DB
	Redis         *redis.Client
	XmrAuthClient *http.Client
}

type SecretSuccessResponse struct {
	Succ    bool
	Message string
	Secret  string
}

var moneroRpcUrl = "http://localhost:28082/json_rpc"

// HANDLERS
// GET
func (inject *Injection) GetIndexView(w http.ResponseWriter, r *http.Request) {
	views.Index().Render(r.Context(), w)
}

func (inject *Injection) GetSignupView(w http.ResponseWriter, r *http.Request) {
	views.Signup().Render(r.Context(), w)
}

func (inject *Injection) GetSigninView(w http.ResponseWriter, r *http.Request) {
	views.Signin().Render(r.Context(), w)
}

func (inject *Injection) GetTransactionPayment(w http.ResponseWriter, r *http.Request) {
	// get transaction data from param id to show in page
	views.Transaction().Render(r.Context(), w)
}

// POST
func (inject *Injection) PostSignup(w http.ResponseWriter, r *http.Request) {

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
			Succ:    false,
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
			Succ:    false,
			Message: "Provide a username",
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	if body.Password == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errMsg := types.JsonResponse{
			Succ:    false,
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
			Succ:    false,
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
			Succ:    false,
			Message: "An error occurred, retry later",
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}
	var queryResult types.User
	userHash, userHashErr := uuid.NewString()
	spoofedUserHash := passwords.Hash256(userHash)
	if userHashErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errMsg := types.JsonResponse{
			Succ:    false,
			Message: "Error, retry later",
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}
	userCheckErr := inject.Psql.QueryRow(
		"INSERT INTO users (hash, username, password, salt) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING RETURNING *;",
		&spoofedUserHash,
		&body.Username,
		hashedPassword,
		newSalt,
	).Scan(&queryResult.Id, &queryResult.Hash, &queryResult.Username, &queryResult.Password, &queryResult.Salt)
	if userCheckErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errMsg := types.JsonResponse{
			Succ:    false,
			Message: "Username not available",
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	succMsg := SecretSuccessResponse{
		Succ:    true,
		Message: "Successfully registered",
		Secret:  userHash,
	}
	json.NewEncoder(w).Encode(succMsg)
}

func (inject *Injection) PostSignin(w http.ResponseWriter, r *http.Request) {
	username, password, _ := r.BasicAuth()
	if username == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errMsg := types.JsonResponse{
			Succ:    false,
			Message: "Provide your username",
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}
	if password == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errMsg := types.JsonResponse{
			Succ:    false,
			Message: "Provide your password",
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	var queryResult types.User
	queryErr := inject.Psql.QueryRow(
		"SELECT * FROM users WHERE username=$1", username).Scan(&queryResult.Id, &queryResult.Hash, &queryResult.Username, &queryResult.Password, &queryResult.Salt)
	if queryErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errMsg := types.JsonResponse{
			Succ:    false,
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
			Succ:    false,
			Message: "Incorrect credentials",
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	ctx := context.Background()
	newSessionId, _ := uuid.NewString()
	hashedSessionId := passwords.Hash256(newSessionId)
	sessionError := inject.Redis.Set(ctx, hashedSessionId, queryResult.Id, 168*time.Hour).Err()
	if sessionError != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errMsg := types.JsonResponse{
			Succ:    false,
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
		Succ:    true,
		Message: "Connected",
	}
	json.NewEncoder(w).Encode(succMsg)
}

func (inject *Injection) PutChangePassword(w http.ResponseWriter, r *http.Request) {
	type changePasswordType struct {
		NewPassword string
		UserHash    string
	}
	var body *changePasswordType
	reqBody, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errMsg := types.JsonResponse{
			Succ:    false,
			Message: "Error, retry later",
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}
	json.Unmarshal(reqBody, &body)

	if body.UserHash == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errMsg := types.JsonResponse{
			Succ:    false,
			Message: "Provide your reset secret",
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	if body.NewPassword == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errMsg := types.JsonResponse{
			Succ:    false,
			Message: "Choose a new password",
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	size, number, upper, special := passwords.VerifyPassword(body.NewPassword)
	if !size || !number || !upper || !special {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errMsg := types.JsonResponse{
			Succ:    false,
			Message: "Password too weak",
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	newHashedPassword, newSalt, hashErr := passwords.HashPassword(body.NewPassword)
	newUserHash, userHashErr := uuid.NewString()
	hashedUserSecret := passwords.Hash256(newUserHash)
	oldUserHash := passwords.Hash256(body.UserHash)

	if hashErr != nil || userHashErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errMsg := types.JsonResponse{
			Succ:    false,
			Message: "Error, retry later",
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	updateUserErr := inject.Psql.QueryRow(
		"UPDATE users SET password=$1, salt=$2, hash=$3 WHERE hash=$4;",
		&newHashedPassword,
		&newSalt,
		&hashedUserSecret,
		&oldUserHash,
	).Err()
	if updateUserErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errMsg := types.JsonResponse{
			Succ:    false,
			Message: "Error, retry later",
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	succMsg := SecretSuccessResponse{
		Succ:    true,
		Message: "Successfully changed password",
		Secret:  newUserHash,
	}
	json.NewEncoder(w).Encode(succMsg)
}

func (inject *Injection) PostBuyerInitTransaction(w http.ResponseWriter, r *http.Request) {
	type transactionBody struct {
		FiatAmount float64
	}
	var body *transactionBody
	reqBody, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errMsg := types.JsonResponse{
			Succ:    false,
			Message: "Error, retry later",
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}
	json.Unmarshal(reqBody, &body)

	if body.FiatAmount == 0 || body.FiatAmount < 10.00 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errMsg := types.JsonResponse{
			Succ:    false,
			Message: "Provide a valid amount in USD",
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	createAccountParams := []byte(`{
		"method":"create_account",
		"params":{"label":"fidexmr-transaction"}
	}`)
	newXmrAccount, accountErr := inject.XmrAuthClient.Post(moneroRpcUrl, "application/json", bytes.NewBuffer(createAccountParams))
	if accountErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errMsg := types.JsonResponse{
			Succ:    false,
			Message: "Error, please retry later",
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	type XmrResult struct {
		AccountIndex uint64
		Address      string
	}
	type XmrCreate struct {
		Id      string
		Jsonrpc string
		Result  XmrResult
	}
	var xmrResp XmrCreate
	xmrResBody, _ := io.ReadAll(newXmrAccount.Body)
	json.Unmarshal(xmrResBody, &xmrResp)
	userId := r.Context().Value("userId")
	var calculatedFee = body.FiatAmount * 0.02
	var buyerDeposit = calculatedFee + body.FiatAmount

	now := time.Now()

	url, err := url.Parse(r.Host)
	if err != nil {
		log.Fatal(err)
	}
	transactionId, trIdErr := uuid.NewString()
	var transactionUrl = fmt.Sprintf("http://%s/transaction?id=%s", url, transactionId) // create url to send to second party

	if trIdErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errMsg := types.JsonResponse{
			Succ:    false,
			Message: "Error, please retry later",
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	queryErr := inject.Psql.QueryRow("INSERT INTO transactions (id, transaction_url, owner_id, transaction_address, fiat_amount, deposit_amount, fees, active, exp_date, deposit_exp_date) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);",
		&transactionId,
		&transactionUrl,
		&userId,
		&xmrResp.Result.Address,
		&body.FiatAmount,
		&buyerDeposit,
		&calculatedFee,
		false,
		now.Add(time.Hour*168), // 7 days for transaction validity
		now.Add(time.Minute*5), // 5 minutes to deposit
	).Err()
	if queryErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errMsg := types.JsonResponse{
			Succ:    false,
			Message: "Error, please retry later",
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	xmrPrice, xmrPriceErr := inject.XmrAuthClient.Get("https://min-api.cryptocompare.com/data/price?fsym=XMR&tsyms=USD,EUR,GBP")
	if xmrPriceErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errMsg := types.JsonResponse{
			Succ:    false,
			Message: "Error, please retry later",
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	fetchXmrPrice, _ := io.ReadAll(xmrPrice.Body)
	var xmrMarketPrices types.XmrMarketPrices
	json.Unmarshal([]byte(fetchXmrPrice), &xmrMarketPrices)
	deposit := moneroapi.FiatToXmrMarketprice(buyerDeposit, xmrMarketPrices.USD)

	uriBodyString := fmt.Sprintf(`{"method": "make_uri", "params":{"address":"%s", "amount":"%s"}}`, xmrResp.Result.Address, deposit)
	createUriParams := []byte(uriBodyString)
	newXmrUri, uriErr := inject.XmrAuthClient.Post(moneroRpcUrl, "application/json", bytes.NewBuffer(createUriParams))
	if uriErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errMsg := types.JsonResponse{
			Succ:    false,
			Message: "Error, please retry later",
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	type FinalResult struct {
		Uri            string
		TransactionUrl string
	}
	type InnerResult struct {
		Uri string
	}
	type XmrUriResult struct {
		Id      int
		Jsonrpc string
		Result  InnerResult
	}
	var xmrLink *XmrUriResult
	fetchUri, _ := io.ReadAll(newXmrUri.Body)
	json.Unmarshal(fetchUri, &xmrLink)
	fmt.Println(xmrLink.Result.Uri)
	var finalResult FinalResult = FinalResult{
		Uri:            xmrLink.Result.Uri,
		TransactionUrl: transactionUrl,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(finalResult)
}

func (inject *Injection) PostSellerContractOk(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "login page")
}

func (inject *Injection) PostBuyerTransactionOk(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "login page")
}

// MIDDLEWARES
func (inject *Injection) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, cookieErr := r.Cookie("fidexmr")
		if cookieErr != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.Background()
		id, getErr := inject.Redis.Get(ctx, cookie.Value).Result()
		if getErr != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), "userId", id))
		next.ServeHTTP(w, r)
	})
}

func (inject *Injection) CheckTransactionExpirationDate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check in db if transaction operation still valid in time
		transactionId := r.URL.Query().Get("id")
		println(transactionId)
		next.ServeHTTP(w, r)
	})
}

func GetRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "bad request", http.StatusMethodNotAllowed)
			return
		}
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
