package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

const kafkaURL = "app-prereq-kafka:9092"

func CreateWatch(c *gin.Context) {
	log.Info("/watch POST watch API")

	produceTopic := "watch"

	watch := WATCH{}
	authHeader := c.Request.Header.Get("Authorization")
	id, err := ParseToken(authHeader)
	if err != nil {
		log.Error("/watch POST watch API error: Unauthorized token error")
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}

	qUser := queryById(id)
	if qUser == nil {
		log.Error("/watch POST watch API error: User not found error")
		c.JSON(http.StatusNotFound, "User not found")
		return
	}

	if c.ShouldBindJSON(&watch) == nil {
		// generate (Version 4) UUID
		uid, _ := uuid.NewRandom()
		watch.ID = uid.String()

		//connect watch to a user by referencing user_id of watch to id of user
		watch.UserId = id

		// get current time in UTC
		// format the time and assign the value to the fields
		watch.WatchCreated = time.Now().UTC().Format("2006-01-02 03:04:05")
		watch.WatchUpdated = watch.WatchCreated
		// for all alerts create proper data
		for i := range watch.Alerts {
			uid_two, _ := uuid.NewRandom()
			watch.Alerts[i].ID = uid_two.String()
			watch.Alerts[i].WatchId = watch.ID
			watch.Alerts[i].AlertCreated = watch.WatchCreated
			watch.Alerts[i].AlertUpdated = watch.WatchCreated
		}
		// add watch to watch table
		if !insertWatch(watch) {
			log.Error("/watch POST watch API error: Db connection error")
			c.JSON(http.StatusBadRequest, "error in watch")
			return
		}
		// add alerts to alert table
		for i := range watch.Alerts {
			//fmt.println("Watch_id")
			//fmt.println(watch.ID)
			if !insertAlert(watch.Alerts[i]) {
				log.Error("/watch POST watch API error: Alerts incorrect")
				c.JSON(http.StatusBadRequest, "Alerts are incorrect")
				return
			}
		}

		resp := watch

		log.Info("/watch POST watch API :SENDING resp to watch topic:\n %s", resp.ID)
		log.Info("/watch POST watch API kafka details sent to :" + kafkaURL + " and topic is: " + produceTopic)
		produce(kafkaURL, produceTopic, resp, "insert")

		// remove watch_id from alerts before sending response
		for i := range resp.Alerts {
			resp.Alerts[i].WatchId = ""
		}

		// RETURN THE INSERTED WATCH
		log.Info("/watch POST watch API succeeded")
		c.JSON(http.StatusCreated, resp)

	} else {
		log.Error("/watch POST watch API error: Fields missing error")
		c.JSON(http.StatusBadRequest, "400 Bad request no queries made")
	}
}

func GetAllWatches(c *gin.Context) {
	log.Info("/watches GET All watches API")

	authHeader := c.Request.Header.Get("Authorization")
	id, err := ParseToken(authHeader)
	if err != nil {
		log.Error("/watches GET All watches API error: Token not found error")
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}
	qUser := queryById(id)
	if qUser == nil {
		log.Error("/watches GET All watches API error: User not found error")
		c.JSON(http.StatusNotFound, "User not found")
		return
	}
	qWatches := queryWatchByUserId(id)
	if qWatches == nil {
		log.Error("/watches GET All watches API error: Watches not found error")
		c.JSON(http.StatusUnauthorized, "Watches not found")
		return
	}
	log.Info("/watches GET All watches API succeeded")
	c.JSON(http.StatusOK, qWatches)
}

func GetWatchById(c *gin.Context) {
	log.Info("/watch Get watch by id API")

	authHeader := c.Request.Header.Get("Authorization")
	id, err := ParseToken(authHeader)
	if err != nil {
		log.Error("/watch Get watch by id API error: token error")
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}

	qUser := queryById(id)
	if qUser == nil {
		log.Error("/watch Get watch by id API error: User not found error")
		c.JSON(http.StatusNotFound, "User not found")
		return
	}
	watch_id := c.Param("id")
	watch := queryByWatchID(watch_id)
	if watch == nil {
		log.Error("/watch Get watch by id API error: Watch id does not exist")
		c.JSON(http.StatusNotFound, "watch with id: "+watch_id+" does not exist")
		return
	}
	if qUser.ID != watch.UserId {
		log.Error("/watch Get watch by id API error: User not owner of the watch error")
		c.JSON(http.StatusUnauthorized, "User not owner of the watch")
		return

	}
	log.Info("/watch Get watch by id API succeeded")
	c.JSON(http.StatusOK, watch)
}

func UpdateWatchById(c *gin.Context) {
	log.Info("/watch PUT watch by id API")

	authHeader := c.Request.Header.Get("Authorization")
	id, err := ParseToken(authHeader)
	if err != nil {
		log.Error("/watch PUT watch by id API error: Token error")
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}

	qUser := queryByID(id)
	if qUser == nil {
		log.Error("/watch PUT watch by id API error: User not found error")
		c.JSON(http.StatusNotFound, "User self not found")
		return
	}
	watch_id := c.Param("id")
	watch := queryByWatchID(watch_id)
	if watch == nil {
		log.Error("/watch PUT watch by id API error: watch id does not exist")
		c.JSON(http.StatusNotFound, "watch with id: "+watch_id+" does not exist")
		return
	}
	if qUser.ID != watch.UserId {
		log.Error("/watch PUT watch by id API error: User not owner of the watch")
		c.JSON(http.StatusUnauthorized, "User not owner of the watch")
		return

	}

	updatedWatch := WATCH{}
	if c.ShouldBindJSON(&updatedWatch) == nil {
		// these values cannot be updated by the user
		updatedWatch.ID = watch.ID
		updatedWatch.UserId = watch.UserId
		updatedWatch.WatchCreated = watch.WatchCreated
		updatedWatch.WatchUpdated = time.Now().UTC().Format("2006-01-02 03:04:05")
		for i := range updatedWatch.Alerts {
			uid_two, _ := uuid.NewRandom()
			updatedWatch.Alerts[i].ID = uid_two.String()
			updatedWatch.Alerts[i].WatchId = updatedWatch.ID
			updatedWatch.Alerts[i].AlertCreated = updatedWatch.WatchUpdated
			updatedWatch.Alerts[i].AlertUpdated = updatedWatch.WatchUpdated
		}
		for i := range watch.Alerts {
			if !deleteAlert(watch.Alerts[i].ID) {
				c.JSON(http.StatusBadRequest, "400 Bad request")
				return
			}
		}

		if !updateWatch(updatedWatch) {
			log.Error("/watch PUT watch by id API error: DB update watch failed")
			c.JSON(http.StatusBadRequest, "400 Bad request")
			return
		}
		for i := range updatedWatch.Alerts {
			if !insertAlert(updatedWatch.Alerts[i]) {
				log.Error("/watch PUT watch by id API error: Db update failed")
				c.JSON(http.StatusBadRequest, "Alerts are incorrect")
				return
			}
		}

		produceTopic := "watch"
		log.Info("SENDING updated watch to watch topic")
		produce(kafkaURL, produceTopic, updatedWatch, "update")
		log.Info("/watch PUT watch by id API succeeded")
		c.Status(http.StatusNoContent)
	} else {
		log.Error("/watch PUT watch by id API error: Fields missing error")
		c.JSON(http.StatusBadRequest, "400 Bad request")
	}
}

func DeleteWatch(c *gin.Context) {
	log.Info("/watch DELETE watch by id API")
	authHeader := c.Request.Header.Get("Authorization")
	id, err := ParseToken(authHeader)
	if err != nil {
		log.Error("/watch DELETE watch by id API error: Token error")
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}

	qUser := queryByID(id)
	if qUser == nil {
		log.Error("/watch DELETE watch by id API error:User not found error")
		c.JSON(http.StatusNotFound, "User self not found")
		return
	}
	watch_id := c.Param("id")
	watch := queryByWatchID(watch_id)
	if watch == nil {
		log.Error("/watch DELETE watch by id API error: Watch id does not exist error")
		c.JSON(http.StatusNotFound, "watch with id: "+watch_id+" does not exist")
		return
	}
	if qUser.ID != watch.UserId {
		log.Error("/watch DELETE watch by id API error: User not owner of the watch error")
		c.JSON(http.StatusUnauthorized, "User not owner of the watch")
		return

	}
	for i := range watch.Alerts {
		if !deleteAlert(watch.Alerts[i].ID) {
			log.Error("/watch DELETE watch by id API error: DB connection failed error")
			c.JSON(http.StatusBadRequest, "400 Bad request")
			return
		}
	}
	if !deleteWatch(watch_id) {
		log.Error("/watch DELETE watch by id API error: DB connection failed error")
		c.JSON(http.StatusBadRequest, "400 Bad request")
		return
	}

	produceTopic := "watch"
	log.Info("SENDING watch to delete on watch topic")
	produce(kafkaURL, produceTopic, *watch, "delete")
	log.Info("/watch DELETE watch by id API succeeded")
	c.Status(http.StatusNoContent)
}
