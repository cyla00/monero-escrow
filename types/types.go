package types

import "time"

type User struct {
	id       string
	username string
	password string
}

// transaction = new wallet account creation
type Transaction struct {
	id              string
	wallet_address  string
	balance         float64
	sending_party   string
	receiving_party string
	created_at      time.Time
}
