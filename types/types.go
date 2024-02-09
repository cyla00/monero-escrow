package types

import (
	"time"

	"github.com/fossoreslp/go-uuid-v4"
)

type User struct {
	id       uuid.UUID
	username string
	password string
}

// transaction = new wallet account creation
type Transaction struct {
	id             uuid.UUID
	owner_id       uuid.UUID
	seller_id      uuid.UUID
	wallet_address string
	balance        float64
	fees           float64
	exp_date       time.Time
	created_at     time.Time
}
