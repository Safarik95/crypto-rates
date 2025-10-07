package api

type BinanceClientInterface interface {
	GetRate(currency string) (float64, error)
}
