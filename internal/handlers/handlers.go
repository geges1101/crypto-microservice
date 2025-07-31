package handlers

import (
	"log"
	"net/http"
	"strconv"

	"crypto-microservice/internal/services"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	cryptoService *services.CryptoService
	priceService  *services.PriceService
}

func NewHandler(cryptoService *services.CryptoService, priceService *services.PriceService) *Handler {
	return &Handler{
		cryptoService: cryptoService,
		priceService:  priceService,
	}
}

// AddCurrencyRequest представляет запрос на добавление криптовалюты
type AddCurrencyRequest struct {
	Symbol string `json:"symbol" binding:"required" example:"bitcoin"`
}

// RemoveCurrencyRequest представляет запрос на удаление криптовалюты
type RemoveCurrencyRequest struct {
	Symbol string `json:"symbol" binding:"required" example:"bitcoin"`
}

// PriceRequest представляет запрос на получение цены
type PriceRequest struct {
	Coin      string `json:"coin" binding:"required" example:"bitcoin"`
	Timestamp int64  `json:"timestamp" binding:"required" example:"1736500490"`
}

// PriceResponse представляет ответ с ценой
type PriceResponse struct {
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Timestamp int64   `json:"timestamp"`
}

// @Summary Добавить криптовалюту в список наблюдения
// @Description Добавляет криптовалюту в список наблюдения и начинает автоматический сбор цен через CoinGecko API.
// При добавлении валюты система начнет собирать её цены каждые 30 секунд.
// @Tags currency
// @Accept json
// @Produce json
// @Param request body AddCurrencyRequest true "Данные для добавления криптовалюты"
// @Success 200 {object} map[string]string "Валюта успешно добавлена"
// @Failure 400 {object} map[string]string "Неверный формат запроса"
// @Failure 500 {object} map[string]string "Ошибка сервера"
// @Router /currency/add [post]
func (h *Handler) AddCurrency(c *gin.Context) {
	var req AddCurrencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Валидация символа валюты
	if req.Symbol == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Currency symbol cannot be empty"})
		return
	}

	// Проверяем, что символ содержит только буквы и цифры
	for _, char := range req.Symbol {
		if !((char >= 'a' && char <= 'z') || (char >= '0' && char <= '9')) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Currency symbol can only contain lowercase letters and numbers"})
			return
		}
	}

	if err := h.cryptoService.AddCurrency(req.Symbol); err != nil {
		log.Printf("Failed to add currency %s: %v", req.Symbol, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add currency", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Currency added successfully", "symbol": req.Symbol})
}

// @Summary Удалить криптовалюту из списка наблюдения
// @Description Удаляет криптовалюту из списка наблюдения и останавливает сбор цен
// @Tags currency
// @Accept json
// @Produce json
// @Param request body RemoveCurrencyRequest true "Данные для удаления криптовалюты"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /currency/remove [delete]
func (h *Handler) RemoveCurrency(c *gin.Context) {
	var req RemoveCurrencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.cryptoService.RemoveCurrency(req.Symbol); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove currency"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Currency removed successfully"})
}

// @Summary Получить цену криптовалюты
// @Description Получает цену криптовалюты в указанный момент времени
// @Tags currency
// @Accept json
// @Produce json
// @Param coin query string true "Символ криптовалюты" example(bitcoin)
// @Param timestamp query int true "Временная метка" example(1736500490)
// @Success 200 {object} PriceResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /currency/price [get]
func (h *Handler) GetPrice(c *gin.Context) {
	coin := c.Query("coin")
	if coin == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Coin parameter is required"})
		return
	}

	timestampStr := c.Query("timestamp")
	if timestampStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Timestamp parameter is required"})
		return
	}

	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid timestamp format"})
		return
	}

	price, err := h.priceService.GetPrice(coin, timestamp)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Price not found"})
		return
	}

	response := PriceResponse{
		Symbol:    price.Currency.Symbol,
		Price:     price.Price,
		Timestamp: price.Timestamp,
	}

	c.JSON(http.StatusOK, response)
}
