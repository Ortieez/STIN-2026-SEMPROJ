package main

import (
	"backend/pkg/api"
	"backend/pkg/storage"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// MockExchangeApi implements api.ExchangeApi interface for testing
type MockExchangeApi struct {
	LatestFunc  func(base string) api.ExchangeApiBaseResponse
	StrongFunc  func(base string) api.ExchangeApiBaseResponse
	WeakFunc    func(base string) api.ExchangeApiBaseResponse
	AverageFunc func(base, selected, from, to string) api.ExchangeApiTimeSeriesResponse
}

func (m *MockExchangeApi) GetLatestExchangeNumbers(base string) api.ExchangeApiBaseResponse {
	return m.LatestFunc(base)
}

func (m *MockExchangeApi) GetStrongestCurrencyToBase(base string) api.ExchangeApiBaseResponse {
	return m.StrongFunc(base)
}

func (m *MockExchangeApi) GetWeakestCurrencyToBase(base string) api.ExchangeApiBaseResponse {
	return m.WeakFunc(base)
}

func (m *MockExchangeApi) GetAverageExchangeRateForCurrencies(base, selected, from, to string) api.ExchangeApiTimeSeriesResponse {
	return m.AverageFunc(base, selected, from, to)
}

func TestMain(m *testing.M) {
	os.Setenv("CACHE_FILE_PATH", "test_main_cache.json")
	os.Setenv("USER_SETTINGS_PATH", "test_settings.json")
	os.Setenv("LOG_FILE_PATH", "test_app.log")
	os.Setenv("AUTH_TOKEN", "test-token")
	os.Setenv("LOGIN_USERNAME", "admin")
	os.Setenv("LOGIN_PASSWORD", "pass")

	code := m.Run()

	os.Remove("test_main_cache.json")
	os.Remove("test_settings.json")
	os.Remove("test_app.log")
	os.Exit(code)
}

func TestLoginEndpoint(t *testing.T) {
	mockApi := &MockExchangeApi{}
	store := storage.NewStorage()
	router := setupRouter(mockApi, store)

	w := httptest.NewRecorder()
	loginData := `{"username":"admin","password":"pass"}`
	req := httptest.NewRequest("POST", "/login", strings.NewReader(loginData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["token"] != "test-token" {
		t.Errorf("Expected token test-token, got %s", resp["token"])
	}
}

func TestProtectedEndpointUnauthorized(t *testing.T) {
	mockApi := &MockExchangeApi{}
	store := storage.NewStorage()
	router := setupRouter(mockApi, store)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/latest", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestLatestEndpoint(t *testing.T) {
	mockApi := &MockExchangeApi{
		LatestFunc: func(base string) api.ExchangeApiBaseResponse {
			return api.ExchangeApiBaseResponse{
				Base: base,
				Date: "2024-01-01",
				Rates: map[string]float64{"USD": 1.1},
			}
		},
	}
	store := storage.NewStorage()
	router := setupRouter(mockApi, store)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/latest?base=EUR", nil)
	req.Header.Set("Authorization", "test-token")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestSettingsEndpoints(t *testing.T) {
	mockApi := &MockExchangeApi{}
	store := storage.NewStorage()
	router := setupRouter(mockApi, store)

	// Test POST settings
	w := httptest.NewRecorder()
	settingsData := `{"baseCurrency":"USD","selectedCurrencies":["EUR","CZK"]}`
	req := httptest.NewRequest("POST", "/settings", strings.NewReader(settingsData))
	req.Header.Set("Authorization", "test-token")
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Test GET settings
	w = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/settings", nil)
	req.Header.Set("Authorization", "test-token")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var settings storage.UserSettings
	json.Unmarshal(w.Body.Bytes(), &settings)
	if settings.BaseCurrency != "USD" {
		t.Errorf("Expected BaseCurrency USD, got %s", settings.BaseCurrency)
	}
}
