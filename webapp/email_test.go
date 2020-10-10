package main

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestIsEmailValid(t *testing.T) {
	email := "boranyldrm21@gmail.com"
	assert.Equal(t, IsEmailValid(email), true)

	email = "boranyldrm21"
	assert.Equal(t, IsEmailValid(email), false)

	email = "boranyldrm21@"
	assert.Equal(t, IsEmailValid(email), false)

	email = "yildirim.b@northeastern.edu"
	assert.Equal(t, IsEmailValid(email), true)
}
