package main

import (
	"bytes"
	"encoding/json"
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

	var response map[string]string
	err := json.Unmarshal([]byte(w.Body.String()), &response)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Nil(t, err)
	assert.Equal(t, response["first_name"], "Jane")
	assert.Equal(t, response["last_name"], "Doe")
	assert.Equal(t, response["username"], "jane.doe@example.com")
}

func TestCreateUserFail(t *testing.T) {
	router := SetupRouter()

	w := httptest.NewRecorder()
	var jsonStr = []byte("{\n  \"last_name\": \"Doe\",\n  \"password\": \"skdjfhskdfjhg\",\n  \"username\": \"jane.doe@example.com\"\n}")
	req, _ := http.NewRequest("POST", "/v1/user", bytes.NewBuffer(jsonStr))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
