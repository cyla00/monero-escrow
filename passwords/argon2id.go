package passwords

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/cyla00/monero-escrow/types"
	"golang.org/x/crypto/argon2"
)

func HashPassword(inputPassword string) (string, error) {

	p := types.Argon2Params{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}

	salt := make([]byte, p.SaltLength)
	_, randRrr := rand.Read(salt)
	if randRrr != nil {
		return "", randRrr
	}
	passwordToBytes := argon2.IDKey([]byte(inputPassword), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)
	saltHuman := base64.RawStdEncoding.EncodeToString(salt)
	passwordHuman := base64.RawStdEncoding.EncodeToString(passwordToBytes)
	hashedPassword := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, p.Memory, p.Iterations, p.Parallelism, saltHuman, passwordHuman)
	return hashedPassword, nil
}
