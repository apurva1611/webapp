package main

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

func BcryptAndSalt(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Print(err.Error())
	}
	return string(hash)
}

func VerifyPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(hash))
	if err != nil {
		return false
	}
	return true
}
