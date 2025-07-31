package main

import (
	"log"
	"os"
	"time"

	"crypto-microservice/internal/config"
	"crypto-microservice/internal/database"
	"crypto-microservice/internal/handlers"
	"crypto-microservice/internal/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Crypto Microservice API
// @version 1.0
// @description Микросервис для сбора, хранения и отображения стоимости криптовалют.
// Система автоматически собирает цены через CoinGecko API и сохраняет их в PostgreSQL.
// Поддерживает добавление/удаление валют из списка наблюдения и получение цен по временным меткам.
// @host localhost:8080
// @BasePath /
func main() {
	// Загружаем конфигурацию
	cfg := config.Load()

	// Инициализируем базу данных
	db, err := database.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Создаем сервисы
	cryptoService := services.NewCryptoService(db)
	priceService := services.NewPriceService(db)

	// Создаем обработчики
	handler := handlers.NewHandler(cryptoService, priceService)

	// Настраиваем Gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API роуты
	api := r.Group("/currency")
	{
		api.POST("/add", handler.AddCurrency)
		api.DELETE("/remove", handler.RemoveCurrency)
		api.GET("/price", handler.GetPrice)
	}

	// Запускаем фоновую задачу для сбора цен
	go func() {
		ticker := time.NewTicker(time.Duration(cfg.PriceUpdateInterval) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := priceService.UpdatePrices(); err != nil {
					log.Printf("Error updating prices: %v", err)
				}
			}
		}
	}()

	// Запускаем сервер
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
