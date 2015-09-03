package security

import (
	"code.google.com/p/go.crypto/bcrypt"
	"math/rand"
)

func passwordClear(b []byte) {
	for i := 0; i < len(b); i++ {
		b[i] = 0
	}
}

// PasswordCrypt encrypts a password.
func PasswordCrypt(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

// PasswordCompare compares a plain text password ( []byte ) to an encrypted byte array
func PasswordCompare(hash []byte, password []byte) error {
	defer passwordClear(password)
	return bcrypt.CompareHashAndPassword(hash, password)
}

// RandomPassword generates a random password of a given length
func RandomPassword(l int) string {
	var letters = []rune("abcdefghjkmnpqrstuvwxyzABCDEFGHJKMNPQRSTUVWXYZ23456789")

	b := make([]rune, l)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
