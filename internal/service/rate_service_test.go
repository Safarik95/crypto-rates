package service_test

import (
	"context"
	"crypto-rates/internal/database"
	"crypto-rates/internal/service"
	"crypto-rates/internal/types"
	"crypto-rates/mocks"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestRateService_FetchAndStoreRates_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBinance := mocks.NewMockBinanceClientInterface(ctrl)
	mockRepo := mocks.NewMockRateRepositoryInterface(ctrl)

	mockBinance.EXPECT().GetRate(gomock.Any(), types.Currency("BTC")).Return(50000.0, nil)
	mockBinance.EXPECT().GetRate(gomock.Any(), types.Currency("ETH")).Return(3000.0, nil)
	mockRepo.EXPECT().SaveRate(gomock.Any(), types.Currency("BTC"), 50000.0).Return(nil)
	mockRepo.EXPECT().SaveRate(gomock.Any(), types.Currency("ETH"), 3000.0).Return(nil)

	rateService := service.NewRateService(mockBinance, mockRepo)
	err := rateService.FetchAndStoreRates(context.Background())

	assert.NoError(t, err)
}

func TestRateService_FetchAndStoreRates_BinanceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBinance := mocks.NewMockBinanceClientInterface(ctrl)
	mockBinance.EXPECT().GetRate(gomock.Any(), types.Currency("BTC")).Return(0.0, errors.New("API error"))

	rateService := service.NewRateService(mockBinance, nil)
	err := rateService.FetchAndStoreRates(context.Background())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get BTC rate")
}

func TestRateService_FetchAndStoreRates_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBinance := mocks.NewMockBinanceClientInterface(ctrl)
	mockRepo := mocks.NewMockRateRepositoryInterface(ctrl)

	mockBinance.EXPECT().GetRate(gomock.Any(), types.Currency("BTC")).Return(50000.0, nil)
	mockRepo.EXPECT().SaveRate(gomock.Any(), types.Currency("BTC"), 50000.0).Return(errors.New("DB error"))

	rateService := service.NewRateService(mockBinance, mockRepo)
	err := rateService.FetchAndStoreRates(context.Background())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to save BTC rate")
}

func TestRateService_GetRateInfo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBinance := mocks.NewMockBinanceClientInterface(ctrl)
	mockRepo := mocks.NewMockRateRepositoryInterface(ctrl)

	expectedRateInfo := &database.RateInfo{
		Currency:     "BTC",
		CurrentPrice: 50000.0,
		MinPrice24h:  49000.0,
		MaxPrice24h:  51000.0,
		Change1h:     "+2.0%",
	}

	mockRepo.EXPECT().GetRateInfo(gomock.Any(), types.Currency("BTC")).Return(expectedRateInfo, nil)
	rateService := service.NewRateService(mockBinance, mockRepo)
	result, err := rateService.GetRateInfo(context.Background(), types.Currency("BTC"))

	assert.NoError(t, err)
	assert.Equal(t, expectedRateInfo, result)
}

func TestRateService_GetAllRateInfo_Success(t *testing.T) {
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
	result := rateService.GetAllRateInfo(context.Background())

	assert.Len(t, result, 2)
	assert.Equal(t, btcInfo, result["BTC"])
	assert.Equal(t, ethInfo, result["ETH"])
}

func TestRateService_GetAllRateInfo_PartialError(t *testing.T) {
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

	mockRepo.EXPECT().GetRateInfo(gomock.Any(), types.Currency("BTC")).Return(btcInfo, nil)
	mockRepo.EXPECT().GetRateInfo(gomock.Any(), types.Currency("ETH")).Return(nil, errors.New("DB error"))
	rateService := service.NewRateService(mockBinance, mockRepo)
	result := rateService.GetAllRateInfo(context.Background())

	assert.Len(t, result, 1)
	assert.Equal(t, btcInfo, result["BTC"])
}
