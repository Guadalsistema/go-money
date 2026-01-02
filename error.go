package currency

import "errors"

var (
	// ErrUnsupportedCurrency is returned when an unsupported currency code is used.
	ErrUnsupportedCurrency = errors.New("unsupported currency")
	ErrDifferentCurrency   = errors.New("different currency types are not supported")
)
