package types

type Currency string

const (
	BTC Currency = "BTC"
	ETH Currency = "ETH"
)

func (c Currency) String() string {
	return string(c)
}

/*func (c Currency) IsValid() bool {
	switch c {
	case BTC, ETH:
		return true
	default:
		return false
	}
}*/
