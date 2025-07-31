package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL         string
	PriceUpdateInterval int
}

func Load() *Config {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:password@localhost:5432/crypto_db?sslmode=disable"
	}

	priceUpdateInterval := 30 // по умолчанию обновляем каждые 30 секунд
	if interval := os.Getenv("PRICE_UPDATE_INTERVAL"); interval != "" {
		if parsed, err := strconv.Atoi(interval); err == nil {
			priceUpdateInterval = parsed
		} else {
			log.Printf("Invalid PRICE_UPDATE_INTERVAL value: %s, using default", interval)
		}
	}

	// TODO: добавить валидацию минимального интервала (не менее 10 секунд)
	// TODO: добавить поддержку конфигурационного файла
	// TODO: добавить метрики для мониторинга

	return &Config{
		DatabaseURL:         databaseURL,
		PriceUpdateInterval: priceUpdateInterval,
	}
}
