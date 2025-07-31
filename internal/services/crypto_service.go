package services

import (
	"crypto-microservice/internal/database"
	"log"

	"gorm.io/gorm"
)

type CryptoService struct {
	db *gorm.DB
}

func NewCryptoService(db *gorm.DB) *CryptoService {
	return &CryptoService{db: db}
}

func (s *CryptoService) AddCurrency(symbol string) error {
	// Проверяем, существует ли уже такая валюта в базе
	var existingCurrency database.Currency
	if err := s.db.Where("symbol = ?", symbol).First(&existingCurrency).Error; err == nil {
		// Если валюта существует, просто активируем её (возможно была деактивирована ранее)
		log.Printf("Currency %s already exists, reactivating...", symbol)
		return s.db.Model(&existingCurrency).Update("is_active", true).Error
	}
	
	// Если валюта не существует, создаём новую запись
	log.Printf("Creating new currency: %s", symbol)
	currency := &database.Currency{
		Symbol:   symbol,
		IsActive: true,
	}
	
	return s.db.Create(currency).Error
}

func (s *CryptoService) RemoveCurrency(symbol string) error {
	// Находим валюту
	var currency database.Currency
	if err := s.db.Where("symbol = ?", symbol).First(&currency).Error; err != nil {
		return err
	}

	// Удаляем все цены для этой валюты
	if err := s.db.Where("currency_id = ?", currency.ID).Delete(&database.Price{}).Error; err != nil {
		return err
	}

	// Удаляем саму валюту
	return s.db.Delete(&currency).Error
}

func (s *CryptoService) GetActiveCurrencies() ([]database.Currency, error) {
	var currencies []database.Currency
	err := s.db.Where("is_active = ?", true).Find(&currencies).Error
	return currencies, err
}
