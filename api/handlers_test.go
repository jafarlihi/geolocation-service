package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServeRetrievalRequest(t *testing.T) {
	req, err := http.NewRequest("GET", "/location", nil)
	if err != nil {
		t.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("ipAddress", "8.8.8.8")
	req.URL.RawQuery = q.Encode()

	dsw := &dataServiceWrapper{ds: &dataServiceMock{}}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(dsw.serveRetrievalRequest)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"IPAddress":"8.8.8.8","CountryCode":"","Country":"","City":"","Latitude":50.5,"Longitude":0,"MysteryValue":0}
`
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestServeRetrievalRequestWithNonExistingIPAddress(t *testing.T) {
	req, err := http.NewRequest("GET", "/location", nil)
	if err != nil {
		t.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("ipAddress", "9.9.9.9")
	req.URL.RawQuery = q.Encode()

	dsw := &dataServiceWrapper{ds: &dataServiceMock{}}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(dsw.serveRetrievalRequest)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
