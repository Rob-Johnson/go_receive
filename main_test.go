package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNonPost(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)

	response := httptest.NewRecorder()
	ReceiveHandler(response, request)

	if response.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected %d method, got %d", http.StatusMethodNotAllowed, response.Code)
	}
}

func TestNoBody(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", nil)

	response := httptest.NewRecorder()
	ReceiveHandler(response, request)

	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected %d method, got %d", http.StatusBadRequest, response.Code)
	}
}

func TestInvalidJSON(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", bytes.NewBufferString("{foo"))

	response := httptest.NewRecorder()
	ReceiveHandler(response, request)

	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected %d method, got %d", http.StatusBadRequest, response.Code)
	}
}

func TestMissingJSON(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", bytes.NewBufferString("{foo:bar}"))

	response := httptest.NewRecorder()
	ReceiveHandler(response, request)

	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected %d method, got %d", http.StatusBadRequest, response.Code)
	}
}

func TestValidJSON(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", bytes.NewBufferString(`{"name":"Galaxy Nexus", "price":"3460.00"}`))

	response := httptest.NewRecorder()
	ReceiveHandler(response, request)

	if response.Code == http.StatusBadRequest {
		t.Errorf("Didn't expect 400 status")
	}

}
