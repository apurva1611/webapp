package main

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestCrypto(t *testing.T) {
	password := "123456"
	hash := BcryptAndSalt(password)
	result := VerifyPassword(hash, password)
	assert.Equal(t, result, true)

	password = "AbCd-*1*"
	hash = BcryptAndSalt(password)
	result = VerifyPassword(hash, password)
	assert.Equal(t, result, true)

	password = "AbCd-*1*"
	hash = BcryptAndSalt(password)
	result = VerifyPassword(hash, "AbCd-*1")
	assert.Equal(t, result, false)
}
