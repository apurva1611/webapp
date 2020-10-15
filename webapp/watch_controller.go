package main

import (
	"fmt"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateWatch(c *gin.Context) {
	watch:=WATCH{}
	authHeader := c.Request.Header.Get("Authorization")
	fmt.Printf(authHeader)
	id, err := ParseToken(authHeader)
	if err != nil {
		c.JSON(http.StatusBadRequest, "No Content")
		return
	}

	qUser := queryById(id)
	if qUser == nil {
		c.JSON(http.StatusUnauthorized, "User not found")
		return
	}
    
	if c.ShouldBindJSON(&watch) == nil {
		fmt.Printf(watch.Zipcode)
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
		for i := range  watch.Alerts{
			uid_two, _ := uuid.NewRandom()
			watch.Alerts[i].ID = uid_two.String()
			watch.Alerts[i].WatchId = watch.ID
			fmt.Println("id")
			fmt.Println(watch.Alerts[i].WatchId)
			watch.Alerts[i].AlertCreated = watch.WatchCreated
			watch.Alerts[i].AlertUpdated = watch.WatchCreated
		}
        // add watch to watch table
		if !insertWatch(watch){
			c.JSON(http.StatusBadRequest, "error in watch")
			return
		}
		// add alerts to alert table
		for i := range  watch.Alerts{
			//fmt.println("Watch_id")
			//fmt.println(watch.ID)
			if !insertAlert(watch.Alerts[i]){
				c.JSON(http.StatusBadRequest, "Alerts are incorrect")
				return
			}
		}
		// remove watch_id from alerts before sending response
		resp := watch
		for i := range resp.Alerts{
			 resp.Alerts[i].WatchId = ""
		}
		// RETURN THE INSERTED WATCH
		c.JSON(http.StatusCreated,resp)
	
	} else {
		fmt.Printf("Error")
		c.JSON(http.StatusBadRequest, "400 Bad request no queries made")
	}
}

func GetAllWatches(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	fmt.Printf(authHeader)
	id, err := ParseToken(authHeader)
	if err != nil {
		c.JSON(http.StatusBadRequest, "unauthorized")
		return
	}
	qUser := queryById(id)
	if qUser == nil {
		c.JSON(http.StatusUnauthorized, "User not found")
		return
	}
	qWatches:= queryWatchByUserId(id)
	if qWatches == nil {
		c.JSON(http.StatusUnauthorized, "Watches not found")
		return
	}
	c.JSON(http.StatusOK,qWatches)
}

func GetWatchById(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	fmt.Printf(authHeader)
	id, err := ParseToken(authHeader)
	if err != nil {
		c.JSON(http.StatusBadRequest, "unauthorized")
		return
	}

	qUser := queryById(id)
	if qUser == nil {
		c.JSON(http.StatusUnauthorized, "User not found")
		return
	}
	watch_id := c.Param("id")
	watch := queryByWatchID(watch_id)
	if watch == nil {
		c.JSON(http.StatusNotFound, "watch with id: "+watch_id+" does not exist")
		return
	}
	if(qUser.ID!=watch.UserId){
		c.JSON(http.StatusUnauthorized, "User not owner of the watch")
		return

	}

	c.JSON(http.StatusOK,watch)
}

func UpdateWatchById(c *gin.Context) {
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
	watch_id := c.Param("id")
	watch := queryByWatchID(watch_id)
	if watch == nil {
		c.JSON(http.StatusNotFound, "watch with id: "+watch_id+" does not exist")
		return
	}
	if(qUser.ID!=watch.UserId){
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
		for i := range updatedWatch.Alerts{
			uid_two, _ := uuid.NewRandom()
			updatedWatch.Alerts[i].ID = uid_two.String()
			updatedWatch.Alerts[i].WatchId = updatedWatch.ID
			updatedWatch.Alerts[i].AlertCreated = updatedWatch.WatchUpdated
			updatedWatch.Alerts[i].AlertUpdated = updatedWatch.WatchUpdated
		}
		for i := range watch.Alerts{
			if !deleteAlert(watch.Alerts[i].ID) {
				c.JSON(http.StatusBadRequest, "400 Bad request")
				return
			}
		}
		

		if !updateWatch(updatedWatch) {
			c.JSON(http.StatusBadRequest, "400 Bad request")
			return
		}
		fmt.Printf("Error here")
		for i := range updatedWatch.Alerts{
			if !insertAlert(updatedWatch.Alerts[i]){
				c.JSON(http.StatusBadRequest, "Alerts are incorrect")
				return
			}
		}

	c.Status(http.StatusNoContent)
	} else {
		c.JSON(http.StatusBadRequest, "400 Bad request")
	}
}

func DeleteWatch(c *gin.Context) {
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
	watch_id := c.Param("id")
	watch := queryByWatchID(watch_id)
	if watch == nil {
		c.JSON(http.StatusNotFound, "watch with id: "+watch_id+" does not exist")
		return
	}
	if(qUser.ID!=watch.UserId){
		c.JSON(http.StatusUnauthorized, "User not owner of the watch")
		return

	}
	for i := range watch.Alerts{
		if !deleteAlert(watch.Alerts[i].ID) {
			c.JSON(http.StatusBadRequest, "400 Bad request")
			return
		}
	}
	if !deleteWatch(watch_id) {
		c.JSON(http.StatusBadRequest, "400 Bad request")
		return
	}
	c.Status(http.StatusNoContent)
}

