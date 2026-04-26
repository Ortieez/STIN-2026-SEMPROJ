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
type MockExchangeApi struct{}

func (m *MockExchangeApi) GetLatestExchangeNumbers(baseCurrency string) api.ExchangeApiBaseResponse {
	return api.ExchangeApiBaseResponse{
		Base: baseCurrency,
		Date: "2024-01-01",
		Rates: map[string]float64{
			"USD": 1.1,
		},
	}
}

func (m *MockExchangeApi) GetStrongestCurrencyToBase(baseCurrency string) api.ExchangeApiBaseResponse {
	return api.ExchangeApiBaseResponse{
		Base:  baseCurrency,
		Date:  "2024-01-01",
		Rates: map[string]float64{"USD": 1.1},
	}
}

func (m *MockExchangeApi) GetWeakestCurrencyToBase(baseCurrency string) api.ExchangeApiBaseResponse {
	return api.ExchangeApiBaseResponse{
		Base:  baseCurrency,
		Date:  "2024-01-01",
		Rates: map[string]float64{"CZK": 25.0},
	}
}

func (m *MockExchangeApi) GetAverageExchangeRateForCurrencies(baseCurrency string, selectedCurrencies string, from string, to string) api.ExchangeApiTimeSeriesResponse {
	return api.ExchangeApiTimeSeriesResponse{
		Base:      baseCurrency,
		StartDate: from,
		EndDate:   to,
		Rates: map[string]map[string]float64{
			from: {"USD": 1.1},
		},
	}
}

func TestMain(m *testing.M) {
	os.Setenv("CACHE_FILE_PATH", "test_main_cache.json")
	code := m.Run()
	os.Remove("test_main_cache.json")
	os.Exit(code)
}

func TestLatestEndpoint(t *testing.T) {
	mockApi := &MockExchangeApi{}
	router := setupRouter(mockApi)

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
	mockApi := &MockExchangeApi{}
	router := setupRouter(mockApi)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/average?forCurrencies=USD&from=invalid&to=invalid", nil)
	router.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Errorf("Expected status 400 for invalid date, got %d", w.Code)
	}
}
