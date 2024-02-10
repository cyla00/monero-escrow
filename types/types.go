package types

import (
	"time"
)

type User struct {
	Id       string
	Hash     string
	Username string
	Password string
	Salt     string
}

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

type Argon2Params struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

type JsonResponse struct {
	Message string
}
