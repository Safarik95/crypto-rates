package types

type Currency string

const (
	BTC Currency = "BTC"
	ETH Currency = "ETH"
)

func (c Currency) String() string {
	return string(c)
}
