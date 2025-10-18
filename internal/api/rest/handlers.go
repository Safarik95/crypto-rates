package rest

import (
	"crypto-rates/internal/service"
	"crypto-rates/internal/types"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

type Handlers struct {
	service *service.RateService
}

func NewHandlers(service *service.RateService) *Handlers {
	return &Handlers{
		service: service,
	}
}

// GetRates godoc
// @Summary Получить все курсы
// @Tags rates
// @Success 200 {object} RatesResponse
// @Router /rates [get]
func (h *Handlers) GetRates(c *gin.Context) {
	rates := h.service.GetAllRateInfo(c.Request.Context())
	if len(rates) == 0 {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "No rate data available"})

		return
	}

	response := RatesResponse{
		Rates: make([]RateResponse, 0, len(rates)),
	}

	for _, rate := range rates {
		response.Rates = append(response.Rates, ConvertRateInfo(rate))
	}
	c.JSON(http.StatusOK, response)
	zap.L().Debug("REST API: returned all rates", zap.Int("count", len(rates)))
}

// GetRate godoc
// @Summary Получить курс валюты
// @Tags rates
// @Param cryptocurrency path string true "BTC или ETH"
// @Success 200 {object} RateResponse
// @Router /rates/{cryptocurrency} [get]
func (h *Handlers) GetRate(c *gin.Context) {
	currency := strings.ToUpper(c.Param("cryptocurrency"))

	if currency == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Currency parameter required"})

		return
	}

	rateInfo, err := h.service.GetRateInfo(c.Request.Context(), types.Currency(currency))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Currency " + currency + " not found"})

		return
	}
	response := ConvertRateInfo(rateInfo)
	c.JSON(http.StatusOK, response)
	zap.L().Debug("REST API: returned rate for currency", zap.String("currency", currency))
}

func (h *Handlers) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func (h *Handlers) SwaggerJSON(c *gin.Context) {
	c.File("./docs/swagger.json")
}
