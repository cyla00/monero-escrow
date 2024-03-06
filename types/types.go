package types

import (
	"time"
)

// transaction = new wallet account creation
type Transaction struct {
	Id            string
	OwnerId       string
	SellerId      string
	WalletAddress string
	Balance       float64
	Fees          float64
	ExpDate       time.Time
	CreatedAt     time.Time
	Active        bool
}

type JsonResponse struct {
	Succ    bool
	Message string
}

type XmrMarketPrices struct {
	USD float64
	EUR float64
	GBP float64
}
