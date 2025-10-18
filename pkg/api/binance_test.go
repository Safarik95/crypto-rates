package api_test

import (
	"context"
	"crypto-rates/internal/types"
	"crypto-rates/pkg/api"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const (
	testHTTPClientTimeout = 10 * time.Second
	testRequestTimeout    = 10 * time.Second
)

func TestBinanceClient_GetRate_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v3/ticker/price?symbol=BTCUSDT", r.URL.String())
		response := `{"symbol": "BTCUSDT", "price": "50000.00"}`
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
	defer server.Close()

	client := api.NewBinanceClient(server.URL+"/api/v3", testHTTPClientTimeout, testRequestTimeout)
	price, err := client.GetRate(context.Background(), types.Currency("BTC"))
	assert.NoError(t, err)
	assert.Equal(t, 50000.00, price)
}

func TestBinanceClient_GetRate_InvalidResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid json`))
	}))
	defer server.Close()

	client := api.NewBinanceClient(server.URL+"/api/v3", testHTTPClientTimeout, testRequestTimeout)
	price, err := client.GetRate(context.Background(), types.Currency("BTC"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "JSON parsing error")
	assert.Equal(t, 0.0, price)
}

func TestBinanceClient_GetRate_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := api.NewBinanceClient(server.URL+"/api/v3", testHTTPClientTimeout, testRequestTimeout)
	price, err := client.GetRate(context.Background(), types.Currency("BTC"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API returned status: 500")
	assert.Equal(t, 0.0, price)
}

func TestBinanceClient_GetRate_InvalidPrice(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `{"symbol": "BTCUSDT", "price": "invalid"}`
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
	defer server.Close()

	client := api.NewBinanceClient(server.URL+"/api/v3", testHTTPClientTimeout, testRequestTimeout)
	price, err := client.GetRate(context.Background(), types.Currency("BTC"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "price conversion error")
	assert.Equal(t, 0.0, price)
}
