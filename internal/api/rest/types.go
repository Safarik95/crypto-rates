package rest

import "crypto-rates/internal/database"

type RateResponse struct {
	Currency    string  `json:"currency"`
	Price       float64 `json:"price"`
	MinPrice24h float64 `json:"min_price_24h"`
	MaxPrice24h float64 `json:"max_price_24h"`
	Change1h    string  `json:"change_1h"`
}

type RatesResponse struct {
	Rates []RateResponse `json:"rates"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func ConvertRateInfo(info *database.RateInfo) RateResponse {
	return RateResponse{
		Currency:    info.Currency,
		Price:       info.CurrentPrice,
		MinPrice24h: info.MinPrice24h,
		MaxPrice24h: info.MaxPrice24h,
		Change1h:    info.Change1h,
	}
}
