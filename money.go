package currency

// Currency represents a monetary amount with currency and numbers of decimal places.
type Currency struct {
	Amount   int64  `json:"amount"`   // Integer value
	Currency string `json:"currency"` // Code of the currency (e.g., "USD", "EUR")
	decimals int16  `json:"decimals"` // Number of decimal places
}

// Some values for common currencies
const (
	USDDecimals  = 2
	EURDecimals  = 2
	YuanDecimals = 2
	MADDecimals  = 2
	JPYDecimals  = 0
)

func (m Currency) Decimals() int16 {
	return m.decimals
}

// ToFloat converts the Currency amount to a float64 representation.
func (m Currency) ToFloat() float64 {
	if m.decimals <= 0 {
		return float64(m.Amount)
	}
	return float64(m.Amount) / float64(pow10(int(m.decimals)))
}

func pow10(n int) int64 {
	result := int64(1)
	for i := 0; i < n; i++ {
		result *= 10
	}
	return result
}

// FromFloat creates a Currency instance from a float64 amount.
func FromFloat(amount float64, currency Currency) Currency {
	scale := pow10(int(currency.decimals))
	return Currency{
		Amount:   int64(amount * float64(scale)),
		Currency: currency.Currency,
		decimals: currency.decimals,
	}
}
