package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

type UserSettings struct {
	BaseCurrency       string   `json:"baseCurrency"`
	SelectedCurrencies []string `json:"selectedCurrencies"`
}

type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
}

type Storage struct {
	settingsPath string
	logPath      string
	mu           sync.Mutex
}

func NewStorage() *Storage {
	_ = godotenv.Load()
	settingsPath := os.Getenv("USER_SETTINGS_PATH")
	if settingsPath == "" {
		settingsPath = "settings.json"
	}
	logPath := os.Getenv("LOG_FILE_PATH")
	if logPath == "" {
		logPath = "app.log"
	}

	return &Storage{
		settingsPath: settingsPath,
		logPath:      logPath,
	}
}

func (s *Storage) GetSettings() (UserSettings, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.settingsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return UserSettings{BaseCurrency: "EUR", SelectedCurrencies: []string{}}, nil
		}
		return UserSettings{}, err
	}

	var settings UserSettings
	err = json.Unmarshal(data, &settings)
	return settings, err
}

func (s *Storage) SaveSettings(settings UserSettings) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.settingsPath, data, 0644)
}

func (s *Storage) Log(level, message string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
	}

	f, err := os.OpenFile(s.logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening log file: %v\n", err)
		return
	}
	defer f.Close()

	data, _ := json.Marshal(entry)
	f.Write(data)
	f.Write([]byte("\n"))
}
