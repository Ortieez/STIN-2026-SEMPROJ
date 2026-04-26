package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Setup mock environment variables for testing
	os.Setenv("BASE_API_URL", "http://localhost:9999")
	os.Setenv("API_URL_LATEST", "/latest")
	os.Setenv("API_URL_CURRENCIES", "/currencies")
	os.Setenv("CACHE_FILE_PATH", "test_main_cache.json")

	code := m.Run()

	os.Remove("test_main_cache.json")
	os.Exit(code)
}

func TestLatestEndpoint(t *testing.T) {
	// Start a mock server for the external API the backend calls
	externalApiMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"base":"EUR","date":"2024-01-01","rates":{"USD":1.1}}`))
	}))
	defer externalApiMock.Close()

	os.Setenv("BASE_API_URL", externalApiMock.URL)

	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/latest?base=EUR", nil)
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	
	data, ok := response["data"].(map[string]interface{})
	if !ok || data["base"] != "EUR" {
		t.Errorf("Expected base EUR in response, got %v", data)
	}
}

func TestAverageEndpointInvalidDate(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/average?forCurrencies=USD&from=invalid&to=invalid", nil)
	router.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Errorf("Expected status 400 for invalid date, got %d", w.Code)
	}
}
