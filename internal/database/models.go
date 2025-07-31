package database

import (
	"time"
)

type Currency struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Symbol    string    `gorm:"uniqueIndex;not null" json:"symbol"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Price struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	CurrencyID uint      `gorm:"not null" json:"currency_id"`
	Currency   Currency  `gorm:"foreignKey:CurrencyID" json:"currency"`
	Price      float64   `gorm:"not null" json:"price"`
	Timestamp  int64     `gorm:"not null;index" json:"timestamp"`
	CreatedAt  time.Time `json:"created_at"`
}
