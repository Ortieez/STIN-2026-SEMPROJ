package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

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

	if _, ok := res.Rates["USD"]; !ok {
		t.Errorf("Expected USD to be strongest")
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

	if _, ok := res.Rates["CZK"]; !ok {
		t.Errorf("Expected CZK to be weakest")
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
