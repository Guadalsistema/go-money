package currency

import (
	"fmt"
	"sync"
)

var (
	registryMu       = sync.RWMutex{}
	currencyDecimals = map[string]int16{
		"USD": USDDecimals,
		"EUR": EURDecimals,
		"CNY": YuanDecimals,
		"MAD": MADDecimals,
		"JPY": JPYDecimals,
	}
)

// RegisterCurrency adds or updates a currency code with its decimals.
func RegisterCurrency(code string, decimals int16) error {
	if code == "" {
		return ErrInvalidCurrency
	}
	if decimals < 0 {
		return ErrInvalidDecimals
	}
	registryMu.Lock()
	defer registryMu.Unlock()
	currencyDecimals[code] = decimals
	return nil
}

// DecimalsFor returns the decimals for a currency code if known.
func DecimalsFor(code string) (int16, bool) {
	if code == "" {
		return 0, false
	}
	registryMu.RLock()
	defer registryMu.RUnlock()
	decimals, ok := currencyDecimals[code]
	return decimals, ok
}

// FromCode creates a Currency using the registered decimals for the code.
func FromCode(code string, amount int64) (Currency, error) {
	decimals, ok := DecimalsFor(code)
	if !ok {
		return Currency{}, fmt.Errorf("%w: %s", ErrUnsupportedCurrency, code)
	}
	return Currency{
		Amount:   amount,
		Currency: code,
		decimals: decimals,
	}, nil
}
