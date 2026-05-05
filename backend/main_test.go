package main

import (
	"backend/pkg/api"
	"backend/pkg/storage"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
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
	if m.LatestFunc != nil {
		return m.LatestFunc(base)
	}
	return api.ExchangeApiBaseResponse{}
}

func (m *MockExchangeApi) GetStrongestCurrencyToBase(base string) api.ExchangeApiBaseResponse {
	if m.StrongFunc != nil {
		return m.StrongFunc(base)
	}
	return api.ExchangeApiBaseResponse{}
}

func (m *MockExchangeApi) GetWeakestCurrencyToBase(base string) api.ExchangeApiBaseResponse {
	if m.WeakFunc != nil {
		return m.WeakFunc(base)
	}
	return api.ExchangeApiBaseResponse{}
}

func (m *MockExchangeApi) GetAverageExchangeRateForCurrencies(base, selected, from, to string) api.ExchangeApiTimeSeriesResponse {
	if m.AverageFunc != nil {
		return m.AverageFunc(base, selected, from, to)
	}
	return api.ExchangeApiTimeSeriesResponse{}
}

func hashString(s string) string {
	hash := sha256.Sum256([]byte(s))
	return hex.EncodeToString(hash[:])
}

var validToken = hashString("admin:pass")
var expectedUserHash = hashString("admin")
var expectedPassHash = hashString("pass")

func TestMain(m *testing.M) {
	os.Setenv("CACHE_FILE_PATH", "test_main_cache.json")
	os.Setenv("USER_SETTINGS_PATH", "test_settings.json")
	os.Setenv("LOG_FILE_PATH", "test_app.log")
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
	loginData := fmt.Sprintf(`{"username":"%s","password":"%s"}`, expectedUserHash, expectedPassHash)
	req := httptest.NewRequest("POST", "/login", strings.NewReader(loginData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["token"] != validToken {
		t.Errorf("Expected token %s, got %s", validToken, resp["token"])
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
	// write settings to trigger the "base == ''" fallback
	store.SaveSettings(storage.UserSettings{BaseCurrency: "GBP", SelectedCurrencies: []string{"USD"}})
	router := setupRouter(mockApi, store)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/latest", nil) // empty base, will fallback to GBP
	req.Header.Set("Authorization", validToken)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	// Test empty settings fallback
	os.Remove("test_settings.json")
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest("GET", "/latest", nil)
	req2.Header.Set("Authorization", validToken)
	router.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w2.Code)
	}
}

func TestStrongestWeakestEndpoints(t *testing.T) {
	mockApi := &MockExchangeApi{
		StrongFunc: func(base string) api.ExchangeApiBaseResponse {
			return api.ExchangeApiBaseResponse{Base: base}
		},
		WeakFunc: func(base string) api.ExchangeApiBaseResponse {
			return api.ExchangeApiBaseResponse{Base: base}
		},
	}
	store := storage.NewStorage()
	router := setupRouter(mockApi, store)

	endpoints := []string{"/strongest", "/weakest"}
	for _, ep := range endpoints {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", ep+"?base=EUR", nil)
		req.Header.Set("Authorization", validToken)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200 for %s, got %d", ep, w.Code)
		}
		
		// Fallback branches
		w2 := httptest.NewRecorder()
		req2 := httptest.NewRequest("GET", ep, nil)
		req2.Header.Set("Authorization", validToken)
		router.ServeHTTP(w2, req2)
		if w2.Code != http.StatusOK {
			t.Errorf("Expected status 200 for empty base on %s, got %d", ep, w2.Code)
		}
	}
}

func TestAverageEndpoint(t *testing.T) {
	mockApi := &MockExchangeApi{
		AverageFunc: func(base, selected, from, to string) api.ExchangeApiTimeSeriesResponse {
			return api.ExchangeApiTimeSeriesResponse{
				Rates: map[string]map[string]float64{
					"2024-01-01": {"USD": 1.0, "CZK": 20.0},
					"2024-01-02": {"USD": 2.0}, // Missing CZK here
				},
			}
		},
	}
	store := storage.NewStorage()
	store.SaveSettings(storage.UserSettings{BaseCurrency: "EUR", SelectedCurrencies: []string{"USD", "CZK"}})
	router := setupRouter(mockApi, store)

	// Valid request verifying arithmetic mean with missing data
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/average?base=EUR&forCurrencies=USD,CZK&from=2024-01-01&to=2024-01-02", nil)
	req.Header.Set("Authorization", validToken)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var respBody struct {
		Data struct {
			Rates map[string]float64 `json:"rates"`
		} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &respBody)

	// USD average: (1.0 + 2.0) / 2 = 1.5
	if respBody.Data.Rates["USD"] != 1.5 {
		t.Errorf("Expected USD average 1.5, got %f", respBody.Data.Rates["USD"])
	}
	// CZK average: (20.0) / 1 = 20.0 (The missing date should be ignored)
	if respBody.Data.Rates["CZK"] != 20.0 {
		t.Errorf("Expected CZK average 20.0, got %f", respBody.Data.Rates["CZK"])
	}

	// Invalid dates
	wDate := httptest.NewRecorder()
	reqDate := httptest.NewRequest("GET", "/average?from=invalid&to=invalid", nil)
	reqDate.Header.Set("Authorization", validToken)
	router.ServeHTTP(wDate, reqDate)
	if wDate.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid dates, got %d", wDate.Code)
	}

	// Missing selected currencies
	os.Remove("test_settings.json")
	wCur := httptest.NewRecorder()
	reqCur := httptest.NewRequest("GET", "/average?from=2024-01-01&to=2024-01-02", nil)
	reqCur.Header.Set("Authorization", validToken)
	router.ServeHTTP(wCur, reqCur)
	if wCur.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for empty currencies, got %d", wCur.Code)
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
	req.Header.Set("Authorization", validToken)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	// Test POST Invalid JSON
	wErr := httptest.NewRecorder()
	reqErr := httptest.NewRequest("POST", "/settings", strings.NewReader("{invalid"))
	reqErr.Header.Set("Authorization", validToken)
	router.ServeHTTP(wErr, reqErr)
	if wErr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", wErr.Code)
	}

	// Test GET settings
	w = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/settings", nil)
	req.Header.Set("Authorization", validToken)
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