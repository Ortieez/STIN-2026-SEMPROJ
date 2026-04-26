package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func getMinMaxRates(rates map[string]float64) (minKey string, minVal float64, maxKey string, maxVal float64) {
	firstLoop := true

	for key, value := range rates {
		if firstLoop {
			minKey, minVal = key, value
			maxKey, maxVal = key, value
			firstLoop = false
			continue
		}

		// Check for Strongest (Min)
		if value < minVal {
			minVal = value
			minKey = key
		}

		// Check for Weakest (Max)
		if value > maxVal {
			maxVal = value
			maxKey = key
		}
	}
	return
}

type ExchangeApiBaseResponse struct {
	Base  string             `json:"base"`  // Base currency for reference
	Date  string             `json:"date"`  // "2026-03-04"
	Rates map[string]float64 `json:"rates"` // "CZK" : 1.05
}

type ExchangeApiTimeSeriesResponse struct {
	Base      string                        `json:"base"`
	StartDate string                        `json:"start_date"`
	EndDate   string                        `json:"end_date"`
	Rates     map[string]map[string]float64 `json:"rates"`
}

type ExchangeApi interface {
	GetLatestExchangeNumbers(baseCurrency string) ExchangeApiBaseResponse
	GetStrongestCurrencyToBase(baseCurrency string) ExchangeApiBaseResponse
	GetWeakestCurrencyToBase(baseCurrency string) ExchangeApiBaseResponse
	GetAverageExchangeRateForCurrencies(baseCurrency string, selectedCurrencies string, from string, to string) ExchangeApiTimeSeriesResponse
}

type ExchangeApiClient struct {
	BaseUrl           string
	PossibleEndpoints map[string]string
}

func (e ExchangeApiClient) GetStrongestCurrencyToBase(baseCurrency string) ExchangeApiBaseResponse {
	data := e.GetLatestExchangeNumbers(baseCurrency)
	strongestKey, strongestVal, _, _ := getMinMaxRates(data.Rates)

	return ExchangeApiBaseResponse{
		Base: baseCurrency,
		Date: data.Date,
		Rates: map[string]float64{
			strongestKey: strongestVal,
		},
	}
}

func (e ExchangeApiClient) GetWeakestCurrencyToBase(baseCurrency string) ExchangeApiBaseResponse {
	data := e.GetLatestExchangeNumbers(baseCurrency)
	_, _, weakestKey, weakestVal := getMinMaxRates(data.Rates)

	return ExchangeApiBaseResponse{
		Base: baseCurrency,
		Date: data.Date,
		Rates: map[string]float64{
			weakestKey: weakestVal,
		},
	}
}

func (e ExchangeApiClient) GetLatestExchangeNumbers(baseCurrency string) ExchangeApiBaseResponse {

	url := e.BaseUrl + e.PossibleEndpoints["LATEST"] + "?base=" + baseCurrency

	resp, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	res := new(ExchangeApiBaseResponse)

	err = json.NewDecoder(resp.Body).Decode(res)
	if err != nil {
		panic(err)
	}

	return *res
}

func (e ExchangeApiClient) GetAverageExchangeRateForCurrencies(baseCurrency string, selectedCurrencies string, from string, to string) ExchangeApiTimeSeriesResponse {
	url := e.BaseUrl + "/" + from + ".." + to + "?base=" + baseCurrency + "&symbols=" + selectedCurrencies

	fmt.Println(url)

	resp, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	res := new(ExchangeApiTimeSeriesResponse)

	err = json.NewDecoder(resp.Body).Decode(res)
	if err != nil {
		panic(err)
	}

	return *res
}

func NewExchangeApiClient() ExchangeApi {
	_ = godotenv.Load()

	baseUrl := os.Getenv("BASE_API_URL")
	latestEndpoint := os.Getenv("API_URL_LATEST")
	currenciesEndpoint := os.Getenv("API_URL_CURRENCIES")

	possibleEndpoints := map[string]string{}
	possibleEndpoints["LATEST"] = latestEndpoint
	possibleEndpoints["CURRENCIES"] = currenciesEndpoint

	return &ExchangeApiClient{
		BaseUrl:           baseUrl,
		PossibleEndpoints: possibleEndpoints,
	}
}
