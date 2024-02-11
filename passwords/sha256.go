package passwords

import (
	"crypto/sha256"
	"encoding/hex"
)

func Hash256(hashctx string) string {
	test := sha256.Sum256([]byte(hashctx))
	te := hex.EncodeToString(test[:])
	return te
}
