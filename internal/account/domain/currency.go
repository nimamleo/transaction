package domain

type Currency string

const (
	USD Currency = "USD"
	EUR Currency = "EUR"
	GBP Currency = "GBP"
)

func (c Currency) IsValid() bool {
	switch c {
	case USD, EUR, GBP:
		return true
	default:
		return false
	}
}

func (c Currency) Code() uint16 {
	switch c {
	case USD:
		return 840
	case EUR:
		return 978
	case GBP:
		return 826
	default:
		return 0
	}
}

func (c Currency) String() string {
	return string(c)
}
