package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type MockLogger struct {
	Logs []string
}

func (m *MockLogger) Log(level, message string) {
	m.Logs = append(m.Logs, level+": "+message)
}

func TestGetMinMaxRates(t *testing.T) {
	rates := map[string]float64{
		"USD": 1.1,
		"CZK": 25.0,
		"GBP": 0.8,
	}

	minKey, minVal, maxKey, maxVal := getMinMaxRates(rates)

	if minKey != "GBP" || minVal != 0.8 {
		t.Errorf("Expected min GBP:0.8, got %s:%f", minKey, minVal)
	}
	if maxKey != "CZK" || maxVal != 25.0 {
		t.Errorf("Expected max CZK:25.0, got %s:%f", maxKey, maxVal)
	}
}

func TestGetLatestExchangeNumbers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := ExchangeApiBaseResponse{
			Base: "EUR",
			Date: "2024-01-01",
			Rates: map[string]float64{
				"USD": 1.08,
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := &ExchangeApiClient{
		BaseUrl: server.URL,
		PossibleEndpoints: map[string]string{
			"LATEST": "/latest",
		},
	}

	res := client.GetLatestExchangeNumbers("EUR")

	if res.Base != "EUR" {
		t.Errorf("Expected base EUR, got %s", res.Base)
	}
	if res.Rates["USD"] != 1.08 {
		t.Errorf("Expected USD rate 1.08, got %f", res.Rates["USD"])
	}
}

func TestGetStrongestCurrencyToBase(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := ExchangeApiBaseResponse{
			Base: "EUR",
			Date: "2024-01-01",
			Rates: map[string]float64{
				"USD": 1.08,
				"CZK": 24.5,
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := &ExchangeApiClient{
		BaseUrl: server.URL,
		PossibleEndpoints: map[string]string{
			"LATEST": "/latest",
		},
	}

	res := client.GetStrongestCurrencyToBase("EUR")

	if _, ok := res.Rates["CZK"]; !ok {
		t.Errorf("Expected CZK to be strongest (highest value)")
	}
}

func TestGetWeakestCurrencyToBase(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := ExchangeApiBaseResponse{
			Base: "EUR",
			Date: "2024-01-01",
			Rates: map[string]float64{
				"USD": 1.08,
				"CZK": 24.5,
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := &ExchangeApiClient{
		BaseUrl: server.URL,
		PossibleEndpoints: map[string]string{
			"LATEST": "/latest",
		},
	}

	res := client.GetWeakestCurrencyToBase("EUR")

	if _, ok := res.Rates["USD"]; !ok {
		t.Errorf("Expected USD to be weakest (lowest value)")
	}
}

func TestGetAverageExchangeRateForCurrencies(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := ExchangeApiTimeSeriesResponse{
			Base:      "EUR",
			StartDate: "2024-01-01",
			EndDate:   "2024-01-02",
			Rates: map[string]map[string]float64{
				"2024-01-01": {"USD": 1.0},
				"2024-01-02": {"USD": 2.0},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := &ExchangeApiClient{
		BaseUrl: server.URL,
	}

	res := client.GetAverageExchangeRateForCurrencies("EUR", "USD", "2024-01-01", "2024-01-02")

	if len(res.Rates) != 2 {
		t.Errorf("Expected 2 days of rates, got %d", len(res.Rates))
	}
}

func TestExchangeApiErrors(t *testing.T) {
	logger := &MockLogger{}
	client := &ExchangeApiClient{
		BaseUrl: "",
		logger:  logger,
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic on invalid URL")
		}
	}()
	client.GetLatestExchangeNumbers("EUR")
}

func TestExchangeApiJsonErrors(t *testing.T) {
	logger := &MockLogger{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	client := &ExchangeApiClient{
		BaseUrl: server.URL,
		logger:  logger,
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic on invalid JSON")
		}
	}()
	client.GetLatestExchangeNumbers("EUR")
}

func TestAverageApiErrors(t *testing.T) {
	logger := &MockLogger{}
	client := &ExchangeApiClient{
		BaseUrl: "",
		logger:  logger,
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic on invalid URL")
		}
	}()
	client.GetAverageExchangeRateForCurrencies("EUR", "USD", "2024-01-01", "2024-01-02")
}

func TestAverageApiJsonErrors(t *testing.T) {
	logger := &MockLogger{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	client := &ExchangeApiClient{
		BaseUrl: server.URL,
		logger:  logger,
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic on invalid JSON")
		}
	}()
	client.GetAverageExchangeRateForCurrencies("EUR", "USD", "2024-01-01", "2024-01-02")
}

func TestNewExchangeApiClient(t *testing.T) {
	os.Setenv("BASE_API_URL", "http://test.com")
	os.Setenv("API_URL_LATEST", "/latest")
	os.Setenv("API_URL_CURRENCIES", "/curr")

	client := NewExchangeApiClient(nil)
	c, ok := client.(*ExchangeApiClient)
	if !ok {
		t.Fatal("Expected *ExchangeApiClient type")
	}

	if c.BaseUrl != "http://test.com" {
		t.Errorf("Expected BaseUrl http://test.com, got %s", c.BaseUrl)
	}
}