package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {
	createDb()
	createTable()
	defer closeDB()
	router := SetupRouter()
	log.Fatal(router.Run(":8080"))
}

// SetupRouter function gets updates User from db
func SetupRouter() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/v1")
	authorized := v1.Group("/user/self")
	fmt.Println(authorized)
	authorized.Use(AuthMW(secret))
	{
		authorized.PUT("", UpdateUserSelf)
		v1.GET("/healthcheck", healthcheck)
		v1.POST("/user", CreateUser)
		// user/:id includes user/self, so routing is handled in GetUserWithId
		v1.GET("/user/:id", GetUserWithID, AuthMW(secret), GetUserSelf)
	}
	// authorized = v1.Group("/watch/")
	// authorized.Use(AuthMW(secret))
	// {
	// 	authorized.POST("", AuthMW(secret), CreatWatch)
	// }
	// grouping watch apis together
	authorized_two := v1.Group("/watch")
	authorized_two.Use(AuthMW(secret))
	{ // post api for watch
		authorized_two.POST("", CreateWatch)
		authorized_two.PUT("/:id", UpdateWatchById)
		authorized_two.GET("/:id", GetWatchById)
		authorized_two.DELETE("/:id", DeleteWatch)
		v1.GET("/watches", GetAllWatches)
	}

	fmt.Printf("http://localhost:8080")

	return router
}
func healthcheck(c *gin.Context) {
	c.JSON(http.StatusOK, "ok")
}

// GetUserSelf function gets User from db
func GetUserSelf(c *gin.Context) {
	// get Authorization header "Bearer <token>"
	authHeader := c.Request.Header.Get("Authorization")

	id, err := ParseToken(authHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "500, Internal server error")
	}

	user := queryByID(id)

	if user == nil {
		c.JSON(http.StatusNotFound, "User self not found")
		return
	}

	c.JSON(http.StatusOK, *user)
}

// UpdateUserSelf function gets updates User from db
func UpdateUserSelf(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	fmt.Printf(authHeader)
	id, err := ParseToken(authHeader)
	if err != nil {
		c.JSON(http.StatusNoContent, "204, No content")
		return
	}

	qUser := queryByID(id)
	if qUser == nil {
		c.JSON(http.StatusNotFound, "User self not found")
		return
	}

	updatedUser := User{}
	if c.ShouldBindJSON(&updatedUser) == nil {
		// these values cannot be updated by the user
		if updatedUser.AccountCreated != "" || updatedUser.AccountUpdated != "" || updatedUser.ID != "" {
			c.JSON(http.StatusBadRequest, "400 Bad request")
		}

		updatedUser.ID = qUser.ID
		updatedUser.Password = BcryptAndSalt(updatedUser.Password)
		updatedUser.AccountCreated = qUser.AccountCreated
		updatedUser.AccountUpdated = time.Now().UTC().Format("2006-01-02 03:04:05")

		if !updateUser(updatedUser) {
			c.JSON(http.StatusBadRequest, "400 Bad request")
			return
		}

		c.JSON(http.StatusOK, "Self updated successfully")
	} else {
		c.JSON(http.StatusBadRequest, "400 Bad request")
	}
}

// CreateUser function gets User from db
func CreateUser(c *gin.Context) {
	user := User{}

	producetest("kafka:9092", "watch", "key", "myfirstmessage")

	if c.ShouldBindJSON(&user) == nil {
		// if username is not a valid email respond 400
		if !IsEmailValid(user.Username) {
			c.JSON(http.StatusBadRequest, "400 Bad request")
			return
		}

		// if password is not a valid password respond 400
		if !IsPasswordValid(user.Password) {
			c.JSON(http.StatusBadRequest, "400 Bad request")
			return
		}

		// generate (Version 4) UUID
		uid, _ := uuid.NewRandom()
		user.ID = uid.String()

		// bcrypt password, library uses salt internally
		hash := BcryptAndSalt(user.Password)
		user.Password = string(hash)

		// get current time in UTC
		// format the time and assign the value to the fields
		user.AccountCreated = time.Now().UTC().Format("2006-01-02 03:04:05")
		user.AccountUpdated = user.AccountCreated

		if !insertUser(user) {
			c.JSON(http.StatusBadRequest, "here400 Bad request")
			return
		}

		// remove the password from response
		resp := user
		resp.Password = ""

		// create JWT token
		token := CreateToken(user.ID)

		c.JSON(http.StatusOK, gin.H{"user": resp, "token": token})
	} else {
		c.JSON(http.StatusBadRequest, "400 Bad request HERE")
	}
}

// GetUserWithID function gets UserID from db
func GetUserWithID(c *gin.Context) {
	id := c.Param("id")

	// if v1/user/self is called, skip this function and move to auth middleware
	if id == "self" {
		c.Next()
		return
	}

	// prevent calling other handlers AuthMW and GetUserSelf
	c.Abort()

	user := queryByID(id)

	if user == nil {
		c.JSON(http.StatusNotFound, "User with id: "+id+" does not exist")
		return
	}

	c.JSON(http.StatusOK, *user)
}

// CreatWatch function
// func CreatWatch(c *gin.Context) {
// 	authHeader := c.Request.Header.Get("Authorization")
// 	fmt.Printf(authHeader)
// 	id, err := ParseToken(authHeader)
// 	if err != nil {
// 		c.JSON(http.StatusNoContent, "204, No content")
// 		return
// 	}
// 	watch := WATCH{}

// 	if c.ShouldBindJSON(&watch) == nil {
// 		if err := postcode.Validate(watch.Zipcode); err != nil {
// 			c.JSON(http.StatusBadRequest, "400 Bad request")
// 		}
// 		// assign user id
// 		watch.UserId = id

// 		// generate (Version 4) UUID
// 		wid, _ := uuid.NewRandom()
// 		watch.ID = wid.String()

// 		// get current time in UTC
// 		// format the time and assign the value to the fields
// 		watch.WatchCreated = time.Now().UTC().Format("2006-01-02 03:04:05")
// 		watch.WatchUpdated = watch.WatchCreated

// 		// generate (Version 4) UUID
// 		aid, _ := uuid.NewRandom()
// 		watch.Alerts.AlertID = aid.String()

// 		// get current time in UTC
// 		// format the time and assign the value to the fields
// 		watch.Alerts.AlertCreated = time.Now().UTC().Format("2006-01-02 03:04:05")
// 		watch.Alerts.AlertUpdated = watch.Alerts.AlertCreated

// 	} else {
// 		c.JSON(http.StatusBadRequest, "400 Bad request")
// 	}
// }
