package passwords

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"unicode"

	"github.com/cyla00/monero-escrow/types"
	"golang.org/x/crypto/argon2"
)

func HashPassword(inputPassword string) (string, string, error) {

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
		return "", "", randRrr
	}
	passwordToBytes := argon2.IDKey([]byte(inputPassword), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)
	saltHuman := base64.RawStdEncoding.EncodeToString(salt)
	passwordHuman := base64.RawStdEncoding.EncodeToString(passwordToBytes)
	hashedPassword := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, p.Memory, p.Iterations, p.Parallelism, saltHuman, passwordHuman)
	return hashedPassword, saltHuman, nil
}

func CheckPasswords(newPassword string, oldPasswordHash string, oldSalt string) bool {

	p := types.Argon2Params{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}
	saltToBytes, _ := base64.RawStdEncoding.DecodeString(oldSalt)
	passwordToBytes := argon2.IDKey([]byte(newPassword), saltToBytes, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)
	newPasswordToString := base64.RawStdEncoding.EncodeToString(passwordToBytes)
	newHashedPassword := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, p.Memory, p.Iterations, p.Parallelism, oldSalt, newPasswordToString)
	return oldPasswordHash == newHashedPassword
}

func Hash256(hashctx string) string {
	test := sha256.Sum256([]byte(hashctx))
	te := hex.EncodeToString(test[:])
	return te
}

func VerifyPassword(s string) (sevenOrMore, number, upper, special bool) {
	letters := 0
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			number = true
			letters++
		case unicode.IsUpper(c):
			upper = true
			letters++
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
			letters++
		case unicode.IsLetter(c) || c == ' ':
			letters++
		default:
			return false, false, false, false
		}
	}
	sevenOrMore = letters >= 8
	return
}
