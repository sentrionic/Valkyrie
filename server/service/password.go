package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/crypto/scrypt"
)

// hashPassword hashes the given password using bcrypt
func hashPassword(password string) (string, error) {
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	// using recommended cost parameters from - https://godoc.org/golang.org/x/crypto/scrypt
	shash, err := scrypt.Key([]byte(password), salt, 32768, 8, 1, 32)
	if err != nil {
		return "", err
	}

	// return hex-encoded string with salt appended to password
	hashedPW := fmt.Sprintf("%s.%s", hex.EncodeToString(shash), hex.EncodeToString(salt))

	return hashedPW, nil
}

// comparePasswords compares the stored password with the supplied one
func comparePasswords(storedPassword string, suppliedPassword string) (bool, error) {
	pwsalt := strings.Split(storedPassword, ".")

	if len(pwsalt) < 2 {
		return false, fmt.Errorf("did not provide a valid hash")
	}

	// check supplied password salted with hash
	salt, err := hex.DecodeString(pwsalt[1])

	if err != nil {
		return false, fmt.Errorf("unable to verify user password")
	}

	shash, err := scrypt.Key([]byte(suppliedPassword), salt, 32768, 8, 1, 32)

	if err != nil {
		return false, fmt.Errorf("unable to verify user password")
	}

	return hex.EncodeToString(shash) == pwsalt[0], nil
}
