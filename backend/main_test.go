package main

import (
	"backend/pkg/api"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
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
	code := m.Run()
	os.Remove("test_main_cache.json")
	os.Exit(code)
}

func TestLatestEndpointMissingBase(t *testing.T) {
	mockApi := &MockExchangeApi{
		LatestFunc: func(base string) api.ExchangeApiBaseResponse {
			return api.ExchangeApiBaseResponse{Base: base}
		},
	}
	router := setupRouter(mockApi)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/latest", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestStrongestEndpointMissingBase(t *testing.T) {
	mockApi := &MockExchangeApi{
		StrongFunc: func(base string) api.ExchangeApiBaseResponse {
			return api.ExchangeApiBaseResponse{Base: base}
		},
	}
	router := setupRouter(mockApi)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/strongest", nil)
	router.ServeHTTP(w, req)

	assertStatus(t, w.Code, http.StatusOK)
}

func TestWeakestEndpointMissingBase(t *testing.T) {
	mockApi := &MockExchangeApi{
		WeakFunc: func(base string) api.ExchangeApiBaseResponse {
			return api.ExchangeApiBaseResponse{Base: base}
		},
	}
	router := setupRouter(mockApi)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/weakest", nil)
	router.ServeHTTP(w, req)

	assertStatus(t, w.Code, http.StatusOK)
}

func assertStatus(t *testing.T, got, want int) {
	if got != want {
		t.Errorf("Expected status %d, got %d", want, got)
	}
}


func TestStrongestEndpoint(t *testing.T) {
	mockApi := &MockExchangeApi{
		StrongFunc: func(base string) api.ExchangeApiBaseResponse {
			return api.ExchangeApiBaseResponse{
				Base: base,
				Date: "2024-01-01",
				Rates: map[string]float64{"GBP": 0.8},
			}
		},
	}
	router := setupRouter(mockApi)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/strongest?base=EUR", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].(map[string]interface{})
	rates := data["rates"].(map[string]interface{})
	if _, ok := rates["GBP"]; !ok {
		t.Errorf("Expected GBP in rates, got %v", rates)
	}
}

func TestWeakestEndpoint(t *testing.T) {
	mockApi := &MockExchangeApi{
		WeakFunc: func(base string) api.ExchangeApiBaseResponse {
			return api.ExchangeApiBaseResponse{
				Base: base,
				Date: "2024-01-01",
				Rates: map[string]float64{"IDR": 15000.0},
			}
		},
	}
	router := setupRouter(mockApi)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/weakest?base=EUR", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].(map[string]interface{})
	rates := data["rates"].(map[string]interface{})
	if _, ok := rates["IDR"]; !ok {
		t.Errorf("Expected IDR in rates, got %v", rates)
	}
}

func TestAverageEndpoint(t *testing.T) {
	mockApi := &MockExchangeApi{
		AverageFunc: func(base, selected, from, to string) api.ExchangeApiTimeSeriesResponse {
			return api.ExchangeApiTimeSeriesResponse{
				Base:      base,
				StartDate: from,
				EndDate:   to,
				Rates: map[string]map[string]float64{
					"2024-01-01": {"USD": 1.0, "CZK": 24.0},
					"2024-01-02": {"USD": 1.2, "CZK": 26.0},
				},
			}
		},
	}
	router := setupRouter(mockApi)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/average?base=EUR&forCurrencies=USD,CZK&from=2024-01-01&to=2024-01-02", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].(map[string]interface{})
	rates := data["rates"].(map[string]interface{})
	
	// Average USD: (1.0 + 1.2) / 2 = 1.1
	// Average CZK: (24.0 + 26.0) / 2 = 25.0
	if rates["USD"] != 1.1 {
		t.Errorf("Expected average USD 1.1, got %v", rates["USD"])
	}
	if rates["CZK"] != 25.0 {
		t.Errorf("Expected average CZK 25.0, got %v", rates["CZK"])
	}
}

func TestAverageEndpointInvalidDate(t *testing.T) {
	mockApi := &MockExchangeApi{}
	router := setupRouter(mockApi)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/average?forCurrencies=USD&from=invalid&to=invalid", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid date, got %d", w.Code)
	}
}

func TestAverageEndpointMissingCurrencies(t *testing.T) {
	mockApi := &MockExchangeApi{}
	router := setupRouter(mockApi)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/average?from=2024-01-01&to=2024-01-02", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for missing currencies, got %d", w.Code)
	}
}
