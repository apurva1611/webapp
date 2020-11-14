package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)


func main() {
	createDb()
	createTable()
	defer closeDB()

	//prometheus.MustRegister(counter)

	router := SetupRouter()

	log.Fatal(router.Run())
}

// SetupRouter function gets updates User from db
func SetupRouter() *gin.Engine {
	router := gin.Default()

	p := newPrometheus("gin")
	p.Use(router)
	// registry := prometheus.NewRegistry()
	// registry.Register(requestsPerMinute)
	//InfoLogger.Println("Starting the application...")
	v1 := router.Group("/v1")
	authorized := v1.Group("/user/self")
	fmt.Println(authorized)
	authorized.Use(AuthMW(secret))
	{
		authorized.PUT("", UpdateUserSelf)
		v1.GET("/healthcheck", healthCheck)
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

func healthCheck(c *gin.Context) {
	err := dbHealthCheck()
	if err != nil {
		log.Error("DB HEALTHCHECK %s", err.Error())
		c.JSON(http.StatusInternalServerError, "db health check failed.")
		os.Exit(1)
	}

	err = kafkaHealthCheck(kafkaURL)
	if err != nil {
		log.Error("KAFKA HEALTHCHECK %s", err.Error())
		c.JSON(http.StatusInternalServerError, "kafka health check failed.")
		os.Exit(2)
	}
	c.JSON(http.StatusOK, "ok")
}

// GetUserSelf function gets User from db
func GetUserSelf(c *gin.Context) {
	log.Info("/user/self Get User with self API")

	// get Authorization header "Bearer <token>"
	authHeader := c.Request.Header.Get("Authorization")

	id, err := ParseToken(authHeader)
	if err != nil {
		log.Error("/user/self Get User with self API: Incorrect token")
		c.JSON(http.StatusInternalServerError, "500, Internal server error")
	}

	user := queryByID(id)

	if user == nil {
		log.Error("/user/self Get User with self API: Token Incorrect Error")
		c.JSON(http.StatusNotFound, "User self not found")
		return
	}
	log.Info("/user/self Get User with self API Succeeded")
	c.JSON(http.StatusOK, *user)
}

// UpdateUserSelf function gets updates User from db
func UpdateUserSelf(c *gin.Context) {
	log.Info("/user/self PUT User with self API")
	authHeader := c.Request.Header.Get("Authorization")
	id, err := ParseToken(authHeader)
	if err != nil {
		log.Error("/user/self PUT User with self API error: Token not found error")
		c.JSON(http.StatusNoContent, "204, No content")
		return
	}

	qUser := queryByID(id)
	if qUser == nil {
		log.Error("/user/self PUT User with self API error: Token incorrect error")
		c.JSON(http.StatusNotFound, "User self not found")
		return
	}

	updatedUser := User{}
	if c.ShouldBindJSON(&updatedUser) == nil {
		// these values cannot be updated by the user
		if updatedUser.AccountCreated != "" || updatedUser.AccountUpdated != "" || updatedUser.ID != "" {
			log.Error("/user/self PUT User with self API error: Fields missing error")
			c.JSON(http.StatusBadRequest, "400 Bad request")
		}

		updatedUser.ID = qUser.ID
		updatedUser.Password = BcryptAndSalt(updatedUser.Password)
		updatedUser.AccountCreated = qUser.AccountCreated
		updatedUser.AccountUpdated = time.Now().UTC().Format("2006-01-02 03:04:05")

		if !updateUser(updatedUser) {
			log.Error("/user/self PUT User with self API error: DB update failed")
			c.JSON(http.StatusBadRequest, "400 Bad request")
			return
		}
		log.Info("/user/self PUT User with self API succeeded")
		c.JSON(http.StatusOK, "Self updated successfully")
	} else {
		log.Error("/user/self PUT User with self API error: Fields missing error")
		c.JSON(http.StatusBadRequest, "400 Bad request")
	}
}

// CreateUser function gets User from db
func CreateUser(c *gin.Context) {
	log.Info("/user Post User Create API")
	user := User{}

	producetest("kafka:9092", "watch", "key", "myfirstmessage")

	if c.ShouldBindJSON(&user) == nil {
		// if username is not a valid email respond 400
		if !IsEmailValid(user.Username) {
			log.Error("/user Post User Create API email validation error")
			c.JSON(http.StatusBadRequest, "400 Bad request")
			return
		}

		// if password is not a valid password respond 400
		if !IsPasswordValid(user.Password) {
			log.Error("/user Post User Create API password validation error")
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
			log.Error("/user Post User Create API error:User already exisit error")
			c.JSON(http.StatusBadRequest, "here400 Bad request")
			return
		}

		// remove the password from response
		resp := user
		resp.Password = ""

		// create JWT token
		token := CreateToken(user.ID)
		log.Info("/user Post User Create API succeeded")
		c.JSON(http.StatusOK, gin.H{"user": resp, "token": token})
	} else {
		log.Error("/user Post User Create API Fields missing error")
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

	log.Info("/user/:id GET User with ID API")

	user := queryByID(id)

	if user == nil {
		log.Error("/user/:id GET User with ID API error: ID does not exisit")
		c.JSON(http.StatusNotFound, "User with id: "+id+" does not exist")
		return
	}
	log.Info("/user/:id GET User with ID API succeeded")
	c.JSON(http.StatusOK, *user)
}

// // NewMetrics creates new Metrics instance.
// func NewMetrics() Metrics {
// 	subsystem := exporter
// 	return Metrics{
// 		TotalScrapes: prometheus.NewCounter(prometheus.CounterOpts{
// 			Namespace: namespace,
// 			Subsystem: subsystem,
// 			Name:      "scrapes_total",
// 			Help:      "Total number of times MySQL was scraped for metrics.",
// 		}),
// 		ScrapeErrors: prometheus.NewCounterVec(prometheus.CounterOpts{
// 			Namespace: namespace,
// 			Subsystem: subsystem,
// 			Name:      "scrape_errors_total",
// 			Help:      "Total number of times an error occurred scraping a MySQL.",
// 		}, []string{"collector"}),
// 		Error: prometheus.NewGauge(prometheus.GaugeOpts{
// 			Namespace: namespace,
// 			Subsystem: subsystem,
// 			Name:      "last_scrape_error",
// 			Help:      "Whether the last scrape of metrics from MySQL resulted in an error (1 for error, 0 for success).",
// 		}),
// 		MySQLUp: prometheus.NewGauge(prometheus.GaugeOpts{
// 			Namespace: namespace,
// 			Name:      "up",
// 			Help:      "Whether the MySQL server is up.",
// 		}),
// 	}
// }

// var (
// 	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
// 		Name: "myapp_processed_ops_total",
// 		Help: "The total number of processed events",
// 	})
// 	counter = prometheus.NewCounter(prometheus.CounterOpts{
// 		Namespace: "logging",
// 		Name:      "my_counter",
// 		Help:      "This is my counter",
// 	})
// )

// func recordMetrics() {
// 	go func() {
// 		for {
// 			opsProcessed.Inc()
// 			time.Sleep(2 * time.Second)
// 		}
// 	}()
// }

// func totalRequest() {
// 	go func() {
// 		for {
// 			counter.Add(rand.Float64() * 5)
// 			time.Sleep(2 * time.Second)
// 		}
// 	}()
// }

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
