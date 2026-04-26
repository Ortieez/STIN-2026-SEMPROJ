package main

import (
	"backend/pkg/api"
	"backend/pkg/cache"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type name interface {
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	exchangeApi := api.NewExchangeApiClient()

	router.Use(cache.Middleware(10 * time.Minute))

	router.GET("/latest", func(c *gin.Context) {
		base := c.Query("base")

		if base == "" {
			base = "EUR"
		}

		latestExchanges := exchangeApi.GetLatestExchangeNumbers(base)

		c.JSON(200, gin.H{"data": latestExchanges})
	})

	router.GET("/strongest", func(c *gin.Context) {
		base := c.Query("base")
		if base == "" {
			base = "EUR"
		}

		strongestExchange := exchangeApi.GetStrongestCurrencyToBase(base)

		c.JSON(200, gin.H{"data": strongestExchange})
	})

	router.GET("/weakest", cache.Middleware(1*time.Hour), func(c *gin.Context) {
		base := c.Query("base")

		if base == "" {
			base = "EUR"
		}

		weakestExchange := exchangeApi.GetWeakestCurrencyToBase(base)

		c.JSON(200, gin.H{"data": weakestExchange})
	})

	router.GET("/average", func(c *gin.Context) {
		base := c.Query("base")
		selectedCurrencies := c.Query("forCurrencies")
		from := c.Query("from")
		to := c.Query("to")

		selectedCurrenciesArr := strings.Split(selectedCurrencies, ",")

		if selectedCurrencies == "" {
			c.JSON(400, gin.H{
				"error": "no selected currencies",
			})
			return
		}

		_, errFrom := time.Parse("2006-01-02", from)
		_, errTo := time.Parse("2006-01-02", to)

		if errFrom != nil || errTo != nil {
			c.JSON(400, gin.H{
				"error": "date format error",
			})
			return
		}

		if base == "" {
			base = "EUR"
		}

		data := exchangeApi.GetAverageExchangeRateForCurrencies(base, selectedCurrencies, from, to)
		averageCurrencies := make(map[string]float64)
		totalEntries := float64(len(data.Rates))

		if totalEntries > 0 {
			for _, innerMap := range data.Rates {
				for _, currency := range selectedCurrenciesArr {
					if val, ok := innerMap[currency]; ok {
						averageCurrencies[currency] += val
					}
				}
			}

			for _, currency := range selectedCurrenciesArr {
				averageCurrencies[currency] /= totalEntries
			}
		}

		res := &api.ExchangeApiBaseResponse{
			Base:  base,
			Date:  fmt.Sprintf("%s..%s", from, to),
			Rates: averageCurrencies,
		}

		c.JSON(200, gin.H{"data": res})
	})

	return router
}

func main() {
	router := setupRouter()
	router.Run("0.0.0.0:3000")
}
