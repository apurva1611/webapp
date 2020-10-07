package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func main() {
	router := SetupRouter()
	router.Run(":8080")
}

func SetupRouter() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/v1")
	{
		v1.POST("/user", CreateUser)
	}

	fmt.Printf("http://localhost:8080")

	return router
}

func CreateUser(c *gin.Context) {
	user := User{}
	if c.ShouldBindJSON(&user) == nil {
		// TODO: check if username already exists

		// if username is not a valid email respond 400
		if !isEmailValid(user.Username) {
			c.String(http.StatusBadRequest, "")
			return
		}

		// generate (Version 4) UUID
		uid, _ := uuid.NewRandom()
		user.ID = uid.String()

		// bcrypt password, library uses salt internally
		hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

		// get current time in UTC
		currentTime := time.Now().UTC()
		// format the time and assign the value to the fields
		user.AccountCreated = currentTime.Format("2006-01-02 03:04:05")
		user.AccountUpdated = user.AccountCreated

		// TODO: add user to the database

		// remove the password from response
		resp := user
		resp.Password = string(hash)

		c.JSON(http.StatusOK, resp)
	} else {
		c.String(http.StatusBadRequest, "")
	}
}
