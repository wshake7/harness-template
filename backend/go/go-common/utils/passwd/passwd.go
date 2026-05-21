package passwd

import (
	"golang.org/x/crypto/bcrypt"
)

func Match(rawPassword, encodedPassword string) bool {
	if encodedPassword == "" {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(encodedPassword), []byte(rawPassword))
	if err == nil {
		return true
	}
	return false
}

func Encode(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(bytes), err
}
