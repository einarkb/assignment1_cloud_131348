package main

import "testing"
import "net/http"
import "net/http/httptest"

// TestFunc is a test for handlerFunc
func TestFunc(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(handlerFunc))
	defer testServer.Close()

	resp, err := http.Get(testServer.URL + "/projectinfo/v1/github.com/apache/kafka")
	if err != nil {
		t.Errorf("Error creating the GET request, %s", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected Statuscode %d, received %d", http.StatusOK, resp.StatusCode)
		return
	}
}
