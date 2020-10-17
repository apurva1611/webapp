package main

type WATCH struct {
	ID           string  `json:"watch_id"`
	UserId       string  `json:"user_id"`
	Zipcode      string  `json:"zipcode" binding:"required"`
	Alerts       []ALERT `json:"alerts" binding:"required"`
	WatchCreated string  `json:"watch_created"`
	WatchUpdated string  `json:"watch_updated"`
}
