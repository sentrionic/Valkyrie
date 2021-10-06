package fixture

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var numberRunes = []rune("1234567890")

// RandStringRunes returns a random latin string of the given length
func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}

func randNumberRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = numberRunes[rand.Intn(len(numberRunes))]
	}

	return string(b)
}

func randStringLowerRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes)/2)]
	}

	return string(b)
}

// RandInt returns a random int within the given range
func RandInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}

// RandID returns a 15 character long numeric string
func RandID() string {
	return randNumberRunes(15)
}

// Username returns a random string that's 4 to 15 characters long
func Username() string {
	return RandStringRunes(RandInt(4, 15))
}

// Email returns a random email ending with @example.com
func Email() string {
	email := fmt.Sprintf("%s@example.com", randStringLowerRunes(RandInt(5, 10)))
	return strings.ToLower(email)
}

// RandStr returns a random string that has the given length
func RandStr(n int) string {
	return RandStringRunes(n)
}

// generateAvatar returns an gravatar using the md5 hash of the email
func generateAvatar(email string) string {
	hash := md5.Sum([]byte(email))
	return fmt.Sprintf("https://gravatar.com/avatar/%s?d=identicon", hex.EncodeToString(hash[:]))
}
