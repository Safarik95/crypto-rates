package database

type RateRepositoryInterface interface {
	SaveRate(currency string, price float64) error
	GetCurrentPrice(currency string) (float64, error)
	GetSimpleDailyStats(currency string) (minPrice, maxPrice float64, err error)
	GetHourlyChangePercent(currency string) string
	GetRateInfo(currency string) (*RateInfo, error)
}
