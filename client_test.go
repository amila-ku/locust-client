package locust

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	locusturl = "http://localhost:8089"

	locustStatsResponce = `{
		"current_response_time_percentile_50": 11, 
		"current_response_time_percentile_95": 22, 
		"errors": [], 
		"fail_ratio": 0.31311475409836065, 
		"state": "running", 
		"stats": [], 
		"total_rps": 9.9, 
		"user_count": 5
	}`

	locustTestStoppedResponce = `{
		"message": "Test stopped", 
		"success": true
	}`

	locustTestStartedResponce = `{
		"message": "Swarming started", 
		"success": true
	}`
)

func TestNewClientURLSetting(t *testing.T) {
	c, err := New(locusturl)
	assert.Nil(t, err)
	url := c.BaseURL.String()
	assert.Equal(t, locusturl, url)
}

func TestGenerateLoad(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		fmt.Fprint(w, locustTestStartedResponce)
	}))

	// Close the server when test finishes.
	defer server.Close()

	c, err := New(server.URL)
	assert.Nil(t, err)
	s, err := c.GenerateLoad(5, 1)
	assert.Nil(t, err)
	assert.Equal(t, "Swarming started", s.Message)
	assert.Equal(t, true, s.Success)
}

func TestStopLoad(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		fmt.Fprint(w, locustTestStoppedResponce)
	}))

	// Close the server when test finishes.
	defer server.Close()

	c, err := New(server.URL)
	assert.Nil(t, err)
	s, err := c.StopLoad()
	assert.Nil(t, err)
	assert.Equal(t, "Test stopped", s.Message)
	assert.Equal(t, true, s.Success)
}

func TestGetStatus(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		fmt.Fprint(w, locustStatsResponce)
	}))

	// Close the server when test finishes.
	defer server.Close()

	client, err := New(server.URL)
	assert.Nil(t, err)
	s, err := client.Stats()
	assert.Nil(t, err)
	assert.Equal(t, 5, s.UserCount)
	assert.Equal(t, "running", s.State)
}
