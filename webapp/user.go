package main

type User struct {
	ID             string `json:"id"`
	FirstName      string `json:"first_name" binding:"required"`
	LastName       string `json:"last_name" binding:"required"`
	Password       string `json:"password,omitempty" binding:"required"`
	Username       string `json:"username" binding:"required"`
	AccountCreated string `json:"account_created"`
	AccountUpdated string `json:"account_updated"`
}
