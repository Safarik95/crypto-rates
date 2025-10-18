package rest_test

import (
	"crypto-rates/internal/api/rest"
	"crypto-rates/internal/database"
	"crypto-rates/internal/service"
	"crypto-rates/internal/types"
	"crypto-rates/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlers_GetRates(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBinance := mocks.NewMockBinanceClientInterface(ctrl)
	mockRepo := mocks.NewMockRateRepositoryInterface(ctrl)
	
	btcInfo := &database.RateInfo{
		Currency:     "BTC",
		CurrentPrice: 50000.0,
		MinPrice24h:  49000.0,
		MaxPrice24h:  51000.0,
		Change1h:     "+2.0%",
	}
	ethInfo := &database.RateInfo{
		Currency:     "ETH",
		CurrentPrice: 3000.0,
		MinPrice24h:  2900.0,
		MaxPrice24h:  3100.0,
		Change1h:     "+1.5%",
	}

	mockRepo.EXPECT().GetRateInfo(gomock.Any(), types.Currency("BTC")).Return(btcInfo, nil)
	mockRepo.EXPECT().GetRateInfo(gomock.Any(), types.Currency("ETH")).Return(ethInfo, nil)

	rateService := service.NewRateService(mockBinance, mockRepo)
	handlers := rest.NewHandlers(rateService)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/rates", handlers.GetRates)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/rates", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}
