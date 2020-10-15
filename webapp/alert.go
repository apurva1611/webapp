package main

type ALERT struct {
	ID             string `json:"alert_id"`
	WatchId        string `json:"watch_id,omitempty"`
	FieldType      string `json:"field_type" binding:"required"`
	Operator       string `json:"operator" binding:"required"`
	Value          int `json:"value" binding:"required"`
	AlertCreated string `json:"alert_created"`
	AlertUpdated string `json:"alert_updated"`
}
