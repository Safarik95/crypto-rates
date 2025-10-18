package api

import (
	"context"
	"crypto-rates/internal/types"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type BinanceClient struct {
	baseURL string
	client  *http.Client
	timeout time.Duration
}

type BinanceTickerResponse struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

func NewBinanceClient(baseURL string, clientTimeout, requestTimeout time.Duration) *BinanceClient {
	return &BinanceClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: clientTimeout},
		timeout: requestTimeout,
	}
}

func (c *BinanceClient) GetRate(ctx context.Context, currency types.Currency) (float64, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	url := fmt.Sprintf("%s/ticker/price?symbol=%sUSDT", c.baseURL, currency.String())

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("HTTP request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	var ticker BinanceTickerResponse
	if err := json.NewDecoder(resp.Body).Decode(&ticker); err != nil {
		return 0, fmt.Errorf("JSON parsing error: %w", err)
	}

	price, err := strconv.ParseFloat(ticker.Price, 64)
	if err != nil {

		return 0, fmt.Errorf("price conversion error: %w", err)
	}

	return price, nil
}
