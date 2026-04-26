package main

import (
	"backend/pkg/api"
	"backend/pkg/auth"
	"backend/pkg/cache"
	"backend/pkg/storage"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type name interface {
}

func setupRouter(exchangeApi api.ExchangeApi, store *storage.Storage) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Public endpoints
	router.POST("/login", auth.LoginHandler)

	// Protected group
	protected := router.Group("/")
	protected.Use(auth.Middleware())
	protected.Use(cache.Middleware(10 * time.Minute))

	protected.GET("/latest", func(c *gin.Context) {
		base := c.Query("base")

		if base == "" {
			settings, _ := store.GetSettings()
			base = settings.BaseCurrency
		}
		if base == "" {
			base = "EUR"
		}

		latestExchanges := exchangeApi.GetLatestExchangeNumbers(base)

		c.JSON(200, gin.H{"data": latestExchanges})
	})

	protected.GET("/strongest", func(c *gin.Context) {
		base := c.Query("base")
		if base == "" {
			settings, _ := store.GetSettings()
			base = settings.BaseCurrency
		}
		if base == "" {
			base = "EUR"
		}

		strongestExchange := exchangeApi.GetStrongestCurrencyToBase(base)

		c.JSON(200, gin.H{"data": strongestExchange})
	})

	protected.GET("/weakest", func(c *gin.Context) {
		base := c.Query("base")
		if base == "" {
			settings, _ := store.GetSettings()
			base = settings.BaseCurrency
		}
		if base == "" {
			base = "EUR"
		}

		weakestExchange := exchangeApi.GetWeakestCurrencyToBase(base)

		c.JSON(200, gin.H{"data": weakestExchange})
	})

	protected.GET("/average", func(c *gin.Context) {
		base := c.Query("base")
		selectedCurrencies := c.Query("forCurrencies")
		from := c.Query("from")
		to := c.Query("to")

		if base == "" {
			settings, _ := store.GetSettings()
			base = settings.BaseCurrency
		}
		if base == "" {
			base = "EUR"
		}

		if selectedCurrencies == "" {
			settings, _ := store.GetSettings()
			selectedCurrencies = strings.Join(settings.SelectedCurrencies, ",")
		}

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

	protected.GET("/settings", func(c *gin.Context) {
		settings, err := store.GetSettings()
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to load settings"})
			return
		}
		c.JSON(200, settings)
	})

	protected.POST("/settings", func(c *gin.Context) {
		var settings storage.UserSettings
		if err := c.ShouldBindJSON(&settings); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}
		if err := store.SaveSettings(settings); err != nil {
			c.JSON(500, gin.H{"error": "Failed to save settings"})
			return
		}
		c.JSON(200, gin.H{"message": "Settings saved"})
	})

	return router
}

func main() {
	store := storage.NewStorage()
	exchangeApi := api.NewExchangeApiClient(store)
	router := setupRouter(exchangeApi, store)
	router.Run("0.0.0.0:3000")
}
