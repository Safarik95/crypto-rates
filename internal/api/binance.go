package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type BinanceClient struct {
	baseURL string
	client  *http.Client
}

type BinanceTickerResponse struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

func NewBinanceClient(baseURL string) *BinanceClient {
	return &BinanceClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *BinanceClient) GetRate(currency string) (float64, error) {
	url := fmt.Sprintf("%s/ticker/price?symbol=%sUSDT", c.baseURL, currency)

	resp, err := c.client.Get(url)
	if err != nil {
		return 0, fmt.Errorf("ошибка HTTP запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API вернуло статус: %d", resp.StatusCode)
	}

	var ticker BinanceTickerResponse
	if err := json.NewDecoder(resp.Body).Decode(&ticker); err != nil {
		return 0, fmt.Errorf("ошибка парсинга JSON: %w", err)
	}

	price, err := strconv.ParseFloat(ticker.Price, 64)
	if err != nil {
		return 0, fmt.Errorf("ошибка конвертации цены: %w", err)
	}

	return price, nil
}
