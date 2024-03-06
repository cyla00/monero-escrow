package routes

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	moneroapi "github.com/cyla00/monero-escrow/monero-api"
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

func (inject *Injection) GetTransactionPayment(w http.ResponseWriter, r *http.Request) {
	// get transaction data from param id to show in page
	views.Transaction().Render(r.Context(), w)
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
		Account_Index uint32
		Address       string
	}
	type XmrCreate struct {
		Id      uint32
		Jsonrpc string
		Result  XmrResult
	}
	var xmrResp *XmrCreate
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

	queryErr := inject.Psql.QueryRow("INSERT INTO transactions (id, account_index, transaction_url, owner_id, transaction_address, fiat_amount, deposit_amount, fees, active, exp_date) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);",
		&transactionId,
		&xmrResp.Result.Account_Index,
		&transactionUrl,
		&userId,
		&xmrResp.Result.Address,
		&body.FiatAmount,
		&buyerDeposit,
		&calculatedFee,
		false,
		now.Add(time.Hour*4), // 4 hours for transaction validity
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
	var finalResult FinalResult = FinalResult{
		Uri:            xmrLink.Result.Uri,
		TransactionUrl: transactionUrl,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(finalResult)
}

func (inject *Injection) PutSellerContractOk(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "login page")
}

func (inject *Injection) PostBuyerTransactionOk(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "login page")
}

// MIDDLEWARES
func (inject *Injection) CheckTransactionExpTime(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		transactionId := r.URL.Query().Get("id")

		var expDate time.Time
		dbErr := inject.Psql.QueryRow("SELECT exp_date FROM transactions WHERE id=$1", &transactionId).Scan(&expDate)
		if dbErr != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		checkDate := expDate.Before(time.Now())
		if checkDate {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
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
