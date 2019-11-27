package crypto

var Path = []string{"commons", "crypto", "crypto.go"}

var Content = `package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
	"{{ .AppRepository }}/commons/log"
)

func SetPassword(password string) string {
	var err error

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error.Println(err)
	}

	return string(hash)
}

func CheckPassword(password string, hashedPassword string) bool {
	var err error

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return err == nil
}

func RandString(size int) string {
	var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	rand.Seed(time.Now().UnixNano())

	b := make([]rune, size)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

func EncryptText(clearText string, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(clearText))

	return hex.EncodeToString(mac.Sum(nil))
}`
