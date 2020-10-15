package main

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// BcryptAndSalt represents functin to create crypted password.
func BcryptAndSalt(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Print(err.Error())
	}
	return string(hash)
}

// VerifyPassword represents functin to varify password.
func VerifyPassword(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}
	return true
}
