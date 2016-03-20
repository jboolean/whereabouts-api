package biz

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"log"
)

const TOKEN_LENGTH = 40

func hash(in string) string {
	h := sha512.New()
	h.Write([]byte(in))

	return base64.URLEncoding.EncodeToString(h.Sum(make([]byte, 0)))[:TOKEN_LENGTH]
}

func makeToken() string {
	b := make([]byte, TOKEN_LENGTH)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return base64.URLEncoding.EncodeToString(b)[:TOKEN_LENGTH]
}

func createHashedPassword(password string) (salt string, hashedPassword string) {
	salt = makeToken()
	hashedPassword = hash(salt + password)
	return
}

func checkPassword(candidate, salt, hashedPassword string) bool {
	return hash(salt+candidate) == hashedPassword
}
