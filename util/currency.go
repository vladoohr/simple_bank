package util

const (
	USD = "USD"
	EUR = "EUR"
	CAN = "CAN"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAN:
		return true
	}

	return false
}
