package main

import (
	"backend/pkg/api"
	"backend/pkg/auth"
	"backend/pkg/cache"
	"backend/pkg/i18n"
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
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, Accept-Language")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Public endpoints
	router.POST("/login", auth.LoginHandler)

	// Protected group (No Cache - for Settings)
	protectedNoCache := router.Group("/")
	protectedNoCache.Use(auth.Middleware())

	protectedNoCache.GET("/settings", func(c *gin.Context) {
		settings, err := store.GetSettings()
		if err != nil {
			c.JSON(500, gin.H{"error": i18n.T(c, "failed_load_settings")})
			return
		}
		c.JSON(200, settings)
	})

	protectedNoCache.POST("/settings", func(c *gin.Context) {
		var settings storage.UserSettings
		if err := c.ShouldBindJSON(&settings); err != nil {
			c.JSON(400, gin.H{"error": i18n.T(c, "invalid_request")})
			return
		}
		if err := store.SaveSettings(settings); err != nil {
			c.JSON(500, gin.H{"error": i18n.T(c, "failed_save_settings")})
			return
		}
		c.JSON(200, gin.H{"message": i18n.T(c, "settings_saved")})
	})

	// Protected group (With Cache - for Data)
	protectedCached := router.Group("/")
	protectedCached.Use(auth.Middleware())
	protectedCached.Use(cache.Middleware(10 * time.Minute))

	protectedCached.GET("/latest", func(c *gin.Context) {
		settings, _ := store.GetSettings()
		base := c.Query("base")
		if base == "" {
			base = settings.BaseCurrency
		}
		if base == "" {
			base = "EUR"
		}

		latestExchanges := exchangeApi.GetLatestExchangeNumbers(base)

		// Filter by selected currencies
		if len(settings.SelectedCurrencies) > 0 {
			filteredRates := make(map[string]float64)
			for _, curr := range settings.SelectedCurrencies {
				if val, ok := latestExchanges.Rates[curr]; ok {
					filteredRates[curr] = val
				}
			}
			latestExchanges.Rates = filteredRates
		}

		c.JSON(200, gin.H{"data": latestExchanges})
	})

	protectedCached.GET("/strongest", func(c *gin.Context) {
		settings, _ := store.GetSettings()
		base := c.Query("base")
		if base == "" {
			base = settings.BaseCurrency
		}
		if base == "" {
			base = "EUR"
		}

		data := exchangeApi.GetLatestExchangeNumbers(base)
		
		// Filter first if currencies are selected
		ratesToCompare := data.Rates
		if len(settings.SelectedCurrencies) > 0 {
			filtered := make(map[string]float64)
			for _, curr := range settings.SelectedCurrencies {
				if val, ok := data.Rates[curr]; ok {
					filtered[curr] = val
				}
			}
			ratesToCompare = filtered
		}

		// Find strongest among filtered
		strongestKey := ""
		var strongestVal float64
		first := true
		for k, v := range ratesToCompare {
			if first || v < strongestVal {
				strongestVal = v
				strongestKey = k
				first = false
			}
		}

		res := api.ExchangeApiBaseResponse{
			Base: base,
			Date: data.Date,
			Rates: map[string]float64{strongestKey: strongestVal},
		}

		c.JSON(200, gin.H{"data": res})
	})

	protectedCached.GET("/weakest", func(c *gin.Context) {
		settings, _ := store.GetSettings()
		base := c.Query("base")
		if base == "" {
			base = settings.BaseCurrency
		}
		if base == "" {
			base = "EUR"
		}

		data := exchangeApi.GetLatestExchangeNumbers(base)
		
		// Filter first if currencies are selected
		ratesToCompare := data.Rates
		if len(settings.SelectedCurrencies) > 0 {
			filtered := make(map[string]float64)
			for _, curr := range settings.SelectedCurrencies {
				if val, ok := data.Rates[curr]; ok {
					filtered[curr] = val
				}
			}
			ratesToCompare = filtered
		}

		// Find weakest among filtered
		weakestKey := ""
		var weakestVal float64
		for k, v := range ratesToCompare {
			if v > weakestVal {
				weakestVal = v
				weakestKey = k
			}
		}

		res := api.ExchangeApiBaseResponse{
			Base: base,
			Date: data.Date,
			Rates: map[string]float64{weakestKey: weakestVal},
		}

		c.JSON(200, gin.H{"data": res})
	})

	protectedCached.GET("/average", func(c *gin.Context) {
		base := c.Query("base")
		selectedCurrencies := c.Query("forCurrencies")
		from := c.Query("from")
		to := c.Query("to")

		settings, _ := store.GetSettings()
		if base == "" {
			base = settings.BaseCurrency
		}
		if base == "" {
			base = "EUR"
		}

		if selectedCurrencies == "" {
			selectedCurrencies = strings.Join(settings.SelectedCurrencies, ",")
		}

		if selectedCurrencies == "" {
			c.JSON(400, gin.H{"error": i18n.T(c, "no_selected_currencies")})
			return
		}

		selectedCurrenciesArr := strings.Split(selectedCurrencies, ",")
		_, errFrom := time.Parse("2006-01-02", from)
		_, errTo := time.Parse("2006-01-02", to)

		if errFrom != nil || errTo != nil {
			c.JSON(400, gin.H{"error": i18n.T(c, "date_format_error")})
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

	return router
}

func main() {
	store := storage.NewStorage()
	exchangeApi := api.NewExchangeApiClient(store)
	router := setupRouter(exchangeApi, store)
	router.Run("0.0.0.0:3000")
}
