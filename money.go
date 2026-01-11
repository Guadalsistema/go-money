package currency

import (
	"encoding/json"
	"fmt"
)

// Currency represents a monetary amount with currency and numbers of decimal places.
type Currency struct {
	Amount   int64  `json:"amount"`   // Integer value
	Currency string `json:"currency"` // Code of the currency (e.g., "USD", "EUR")
	decimals int16
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

type currencyJSON struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
	Decimals *int16 `json:"decimals,omitempty"`
}

// MarshalJSON ensures decimals are included in the JSON representation.
func (m Currency) MarshalJSON() ([]byte, error) {
	decimals := m.decimals
	return json.Marshal(currencyJSON{
		Amount:   m.Amount,
		Currency: m.Currency,
		Decimals: &decimals,
	})
}

// UnmarshalJSON restores the currency and decimals, inferring decimals when omitted.
func (m *Currency) UnmarshalJSON(data []byte) error {
	var tmp currencyJSON
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	if tmp.Currency == "" {
		return ErrInvalidCurrency
	}
	decimals := tmp.Decimals
	if decimals == nil {
		derived, ok := DecimalsFor(tmp.Currency)
		if !ok {
			return ErrUnsupportedCurrency
		}
		decimals = &derived
	}
	if *decimals < 0 {
		return fmt.Errorf("invalid decimals: %d", *decimals)
	}
	m.Amount = tmp.Amount
	m.Currency = tmp.Currency
	m.decimals = *decimals
	return nil
}

// ToFloat converts the Currency amount to major units as float64.
func (m Currency) ToFloat() float64 {
	if m.decimals <= 0 {
		return float64(m.Amount)
	}
	return float64(m.Amount) / float64(pow10(int(m.decimals)))
}

// ToMinorFloat converts the Currency amount to minor units as float64.
func (m Currency) ToMinorFloat() float64 {
	return float64(m.Amount)
}

// DecimalString renders the Currency amount as a decimal string in major units.
func (m Currency) DecimalString() string {
	if m.decimals <= 0 {
		return fmt.Sprintf("%d", m.Amount)
	}
	sign := ""
	abs := m.Amount
	if abs < 0 {
		sign = "-"
		abs = -abs
	}
	scale := pow10(int(m.decimals))
	intPart := abs / scale
	fracPart := abs % scale
	return fmt.Sprintf("%s%d.%0*d", sign, intPart, int(m.decimals), fracPart)
}

func pow10(n int) int64 {
	result := int64(1)
	for range n {
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
