package services

import (
	"crypto-microservice/internal/database"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type PriceService struct {
	db *gorm.DB
}

type CoinGeckoResponse struct {
	Prices map[string]float64 `json:"prices"`
}

func NewPriceService(db *gorm.DB) *PriceService {
	return &PriceService{db: db}
}

func (s *PriceService) UpdatePrices() error {
	currencies, err := s.GetActiveCurrencies()
	if err != nil {
		log.Printf("Failed to get active currencies: %v", err)
		return err
	}

	if len(currencies) == 0 {
		log.Println("No active currencies found, skipping price update")
		return nil
	}

	log.Printf("Starting price update for %d currencies", len(currencies))
	
	for _, currency := range currencies {
		log.Printf("Fetching price for %s...", currency.Symbol)
		
		price, err := s.fetchPrice(currency.Symbol)
		if err != nil {
			log.Printf("Error fetching price for %s: %v", currency.Symbol, err)
			continue
		}

		priceRecord := &database.Price{
			CurrencyID: currency.ID,
			Price:      price,
			Timestamp:  time.Now().Unix(),
		}

		if err := s.db.Create(priceRecord).Error; err != nil {
			log.Printf("Error saving price for %s: %v", currency.Symbol, err)
		} else {
			log.Printf("Successfully updated price for %s: $%.2f", currency.Symbol, price)
		}
	}

	log.Println("Price update cycle completed")
	return nil
}

func (s *PriceService) fetchPrice(symbol string) (float64, error) {
	// Используем CoinGecko API для получения цен
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd", symbol)
	
	// Добавляем timeout для HTTP запроса
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	resp, err := client.Get(url)
	if err != nil {
		return 0, fmt.Errorf("HTTP request failed for %s: %w", symbol, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API returned status %d for %s", resp.StatusCode, symbol)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body for %s: %w", symbol, err)
	}

	var result map[string]map[string]float64
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, fmt.Errorf("failed to parse JSON response for %s: %w", symbol, err)
	}

	if priceData, exists := result[symbol]; exists {
		if price, exists := priceData["usd"]; exists {
			if price <= 0 {
				return 0, fmt.Errorf("invalid price (%.2f) for %s", price, symbol)
			}
			return price, nil
		}
	}

	return 0, fmt.Errorf("price not found for %s in API response", symbol)
}

func (s *PriceService) GetPrice(symbol string, timestamp int64) (*database.Price, error) {
	var price database.Price

	// Ищем точное время
	err := s.db.Preload("Currency").
		Joins("JOIN currencies ON prices.currency_id = currencies.id").
		Where("currencies.symbol = ? AND prices.timestamp = ?", symbol, timestamp).
		First(&price).Error

	if err == nil {
		return &price, nil
	}

	// Если не найдено, ищем ближайшее время
	err = s.db.Preload("Currency").
		Joins("JOIN currencies ON prices.currency_id = currencies.id").
		Where("currencies.symbol = ?", symbol).
		Order("ABS(prices.timestamp - " + strconv.FormatInt(timestamp, 10) + ")").
		First(&price).Error

	if err != nil {
		return nil, err
	}

	return &price, nil
}

func (s *PriceService) GetActiveCurrencies() ([]database.Currency, error) {
	var currencies []database.Currency
	err := s.db.Where("is_active = ?", true).Find(&currencies).Error
	return currencies, err
}
