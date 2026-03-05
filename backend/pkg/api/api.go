package api

import (
	_ "bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type ExchangeApiBaseResponse struct {
	Base  string             `json:"base"`  // Base currency for reference
	Date  string             `json:"date"`  // "2026-03-04"
	Rates map[string]float64 `json:"rates"` // "CZK" : 1.05
}

type ExchangeApi interface {
	GetLatestExchangeNumbers() ExchangeApiBaseResponse
}

type ExchangeApiClient struct {
	baseUrl           string
	possibleEndpoints map[string]string
}

func (e ExchangeApiClient) GetLatestExchangeNumbers() ExchangeApiBaseResponse {

	fmt.Println("---")
	fmt.Println(e.baseUrl + e.possibleEndpoints["LATEST"])
	fmt.Println("---")

	resp, err := http.Get(e.baseUrl + e.possibleEndpoints["LATEST"])

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

func NewExchangeApiClient() ExchangeApi {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	baseUrl := os.Getenv("BASE_API_URL")
	latestEndpoint := os.Getenv("API_URL_LATEST")
	currenciesEndpoint := os.Getenv("API_URL_CURRENCIES")

	possibleEndpoints := map[string]string{}
	possibleEndpoints["LATEST"] = latestEndpoint
	possibleEndpoints["CURRENCIES"] = currenciesEndpoint

	return &ExchangeApiClient{
		baseUrl:           baseUrl,
		possibleEndpoints: possibleEndpoints,
	}
}
