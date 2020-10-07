package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUserSuccess(t *testing.T) {
	router := SetupRouter()

	w := httptest.NewRecorder()
	var jsonStr = []byte("{\n  \"first_name\": \"Jane\",\n  \"last_name\": \"Doe\",\n  \"password\": \"skdjfhskdfjhg\",\n  \"username\": \"jane.doe@example.com\"\n}")
	req, _ := http.NewRequest("POST", "/v1/user", bytes.NewBuffer(jsonStr))
	router.ServeHTTP(w, req)

	log.Print(w.Body.String())
	var response map[string]interface{}
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, response, "token")
	assert.Contains(t, response, "user")
}

func TestCreateUserFail(t *testing.T) {
	router := SetupRouter()

	w := httptest.NewRecorder()
	var jsonStr = []byte("{\n  \"last_name\": \"Doe\",\n  \"password\": \"skdjfhskdfjhg\",\n  \"username\": \"jane.doe@example.com\"\n}")
	req, _ := http.NewRequest("POST", "/v1/user", bytes.NewBuffer(jsonStr))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetUserSelf(t *testing.T) {
	router := SetupRouter()

	// create user to get the token
	w := httptest.NewRecorder()
	var jsonStr = []byte("{\n  \"first_name\": \"Jane\",\n  \"last_name\": \"Doe\",\n  \"password\": \"skdjfhskdfjhg\",\n  \"username\": \"jane.doe@example.com\"\n}")
	req, _ := http.NewRequest("POST", "/v1/user", bytes.NewBuffer(jsonStr))
	router.ServeHTTP(w, req)

	log.Print(w.Body.String())
	var response map[string]interface{}
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)

	token := response["token"]

	// try to get self
	w = httptest.NewRecorder()
	var bearer string
	if val, ok := token.(string); ok {
		bearer = "Bearer " + val
	}
	log.Print(bearer)
	req, _ = http.NewRequest("GET", "/v1/user/self", nil)
	req.Header.Set("Authorization", bearer)
	router.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK)

	log.Print(w.Body.String())
	var response2 map[string]interface{}
	err = json.Unmarshal([]byte(w.Body.String()), &response2)
	assert.Nil(t, err)

	assert.Contains(t, response2, "id")
}
