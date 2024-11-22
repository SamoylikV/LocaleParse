package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	GoogleCredentialsJSON string
	SpreadsheetID         string
	RuReadRange           string
	EngReadRange          string
	RedisAddr             string
	RedisPassword         string
	RedisDB               string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}
	cfg := &Config{
		GoogleCredentialsJSON: os.Getenv("GOOGLE_CREDENTIALS_JSON"),
		SpreadsheetID:         os.Getenv("SPREADSHEET_ID"),
		RuReadRange:           os.Getenv("RU_READ_RANGE"),
		EngReadRange:          os.Getenv("ENG_READ_RANGE"),
		RedisAddr:             os.Getenv("REDDIS_ADDRESS"),
		RedisPassword:         os.Getenv("REDDIS_PASSWORD"),
		RedisDB:               os.Getenv("REDDIS_DB"),
	}
	if cfg.GoogleCredentialsJSON == "" || cfg.SpreadsheetID == "" || cfg.RuReadRange == "" || cfg.EngReadRange == "" || cfg.RedisAddr == "" || cfg.RedisPassword == "" || cfg.RedisDB == "" {
		return nil, fmt.Errorf("incomplete config: ensure all variables are set")
	}
	return cfg, nil
}