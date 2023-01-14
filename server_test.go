package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// executeRequest, creates a new ResponseRecorder
// then executes the request by calling ServeHTTP in the router
// after which the handler writes the response to the response recorder
// which we can then inspect.
func executeRequest(req *http.Request, s *Server) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.Router.ServeHTTP(rr, req)

	return rr
}

// checkResponseCode is a simple utility to check the response code
// of the response
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestMagazines(t *testing.T) {
	// Create a New Server Struct
	s := CreateNewServer()
	// Mount Handlers
	s.MountHandlers("mongodb://root:password@localhost:27017/tcp")

	// Create a New Request
	req, _ := http.NewRequest("GET", "/", nil)

	// Execute Request
	response := executeRequest(req, s)

	// Check the response code
	checkResponseCode(t, http.StatusOK, response.Code)
	want := "Hello World!"
	got := response.Body.String()

	if got != want {
		t.Errorf("got = %v want = %v", got, want)
	}
}
