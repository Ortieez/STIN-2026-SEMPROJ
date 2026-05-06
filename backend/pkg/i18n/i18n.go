package i18n

import (
	"github.com/gin-gonic/gin"
)

var translations = map[string]map[string]string{
	"en": {
		"unauthorized":           "Unauthorized",
		"invalid_request":        "Invalid request",
		"invalid_credentials":    "Invalid credentials",
		"no_selected_currencies": "No selected currencies",
		"date_format_error":      "Date format error",
		"failed_load_settings":   "Failed to load settings",
		"failed_save_settings":   "Failed to save settings",
		"settings_saved":         "Settings saved",
	},
	"cs": {
		"unauthorized":           "Neautorizovaný přístup",
		"invalid_request":        "Neplatný požadavek",
		"invalid_credentials":    "Neplatné přihlašovací údaje",
		"no_selected_currencies": "Nebyly vybrány žádné měny",
		"date_format_error":      "Chyba formátu data",
		"failed_load_settings":   "Nepodařilo se načíst nastavení",
		"failed_save_settings":   "Nepodařilo se uložit nastavení",
		"settings_saved":         "Nastavení uloženo",
	},
}

func GetLanguage(c *gin.Context) string {
	lang := c.GetHeader("Accept-Language")
	if lang == "cs" || lang == "cs-CZ" {
		return "cs"
	}
	return "en"
}

func T(c *gin.Context, key string) string {
	lang := GetLanguage(c)
	if val, ok := translations[lang][key]; ok {
		return val
	}
	return key
}
