package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Whitebox test
func TestHome(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Home)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMovedPermanently)
	}

	body := strings.TrimSpace(rr.Body.String())
	if len(body) < 1 {
		t.Error("handler returned unexpected body")
	}
}
