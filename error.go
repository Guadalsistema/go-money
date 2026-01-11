package currency

import "errors"

var (
	// ErrUnsupportedCurrency is returned when an unsupported currency code is used.
	ErrUnsupportedCurrency = errors.New("unsupported currency")
	ErrDifferentCurrency   = errors.New("different currency types are not supported")
	ErrInvalidCurrency     = errors.New("invalid currency code")
	ErrInvalidDecimals     = errors.New("invalid currency decimals")
	ErrInvalidAmount       = errors.New("invalid amount")
	ErrAmountOverflow      = errors.New("amount overflows int64")
)
