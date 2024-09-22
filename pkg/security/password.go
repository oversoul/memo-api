package security

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"hash/crc32"
	"math/rand"

	"golang.org/x/crypto/bcrypt"
)

type Token struct {
	Plain string
	Token string
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func HashEquals(a, b string) bool {

	tok := fmt.Sprintf("%x", sha256.Sum256([]byte(b)))

	// Convert the strings to byte slices
	aBytes := []byte(a)
	bBytes := []byte(tok)

	// Compare the lengths first
	if len(aBytes) != len(bBytes) {
		return false
	}

	// Use subtle.ConstantTimeCompare to compare the byte slices
	return subtle.ConstantTimeCompare(aBytes, bBytes) == 1
}

func GenerateTokenString() Token {
	tokenEntropy := randSeq(40)

	plain := fmt.Sprintf(
		"%s%s%s",
		"", // token prefix
		tokenEntropy,
		fmt.Sprintf("%08x", crc32.ChecksumIEEE([]byte(tokenEntropy))),
	)

	token := fmt.Sprintf("%x", sha256.Sum256([]byte(plain)))
	return Token{Plain: plain, Token: token}
}

func randSeq(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func CompareHashToPassword(a, b string) error {
	return bcrypt.CompareHashAndPassword([]byte(a), []byte(b))
}
