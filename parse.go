package currency

import (
	"math"
	"strings"
)

// RoundingMode controls how extra precision is handled when parsing.
type RoundingMode int

const (
	RoundTruncate RoundingMode = iota
	RoundHalfUp
	RoundBankers
)

type parseConfig struct {
	rounding RoundingMode
	decimals *int16
}

// ParseOption configures ParseAmount.
type ParseOption func(*parseConfig)

// WithRounding sets the rounding mode for ParseAmount.
func WithRounding(mode RoundingMode) ParseOption {
	return func(cfg *parseConfig) {
		cfg.rounding = mode
	}
}

// WithDecimals overrides the currency decimals for ParseAmount.
func WithDecimals(decimals int16) ParseOption {
	return func(cfg *parseConfig) {
		cfg.decimals = &decimals
	}
}

// ParseAmount converts a decimal string amount into minor units for a currency.
func ParseAmount(code, raw string, opts ...ParseOption) (Currency, error) {
	cfg := parseConfig{rounding: RoundTruncate}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	decimals, err := resolveDecimals(code, cfg.decimals)
	if err != nil {
		return Currency{}, err
	}
	amount, err := parseDecimal(raw, decimals, cfg.rounding)
	if err != nil {
		return Currency{}, err
	}
	return Currency{
		Amount:   amount,
		Currency: code,
		decimals: decimals,
	}, nil
}

// MustParseAmount panics if ParseAmount fails.
func MustParseAmount(code, raw string, opts ...ParseOption) Currency {
	amount, err := ParseAmount(code, raw, opts...)
	if err != nil {
		panic(err)
	}
	return amount
}

// ParseAmountOrZero returns a zero amount when parsing fails.
func ParseAmountOrZero(code, raw string, opts ...ParseOption) Currency {
	cfg := parseConfig{rounding: RoundTruncate}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	decimals, err := resolveDecimals(code, cfg.decimals)
	if err != nil {
		return Currency{Currency: code}
	}
	amount, err := parseDecimal(raw, decimals, cfg.rounding)
	if err != nil {
		return Currency{
			Amount:   0,
			Currency: code,
			decimals: decimals,
		}
	}
	return Currency{
		Amount:   amount,
		Currency: code,
		decimals: decimals,
	}
}

func resolveDecimals(code string, override *int16) (int16, error) {
	if code == "" {
		return 0, ErrInvalidCurrency
	}
	if override != nil {
		if *override < 0 {
			return 0, ErrInvalidDecimals
		}
		return *override, nil
	}
	decimals, ok := DecimalsFor(code)
	if !ok {
		return 0, ErrUnsupportedCurrency
	}
	return decimals, nil
}

func parseDecimal(raw string, decimals int16, rounding RoundingMode) (int64, error) {
	if decimals < 0 {
		return 0, ErrInvalidDecimals
	}
	input := strings.TrimSpace(raw)
	if input == "" {
		return 0, ErrInvalidAmount
	}
	sign := int64(1)
	if input[0] == '-' {
		sign = -1
		input = input[1:]
	} else if input[0] == '+' {
		input = input[1:]
	}
	if input == "" {
		return 0, ErrInvalidAmount
	}
	dot := strings.IndexByte(input, '.')
	if dot >= 0 && strings.LastIndexByte(input, '.') != dot {
		return 0, ErrInvalidAmount
	}
	intPart := input
	fracPart := ""
	if dot >= 0 {
		intPart = input[:dot]
		fracPart = input[dot+1:]
	}
	if intPart == "" {
		intPart = "0"
	}
	intValue, err := parseDigits(intPart)
	if err != nil {
		return 0, err
	}
	scale := pow10(int(decimals))
	amount, err := safeMul(intValue, scale)
	if err != nil {
		return 0, err
	}
	if fracPart != "" {
		if err := validateDigits(fracPart); err != nil {
			return 0, err
		}
	}
	if decimals == 0 {
		if fracPart != "" {
			if shouldRound(fracPart, rounding, amount) {
				amount, err = safeAdd(amount, 1)
				if err != nil {
					return 0, err
				}
			}
		}
		return amount * sign, nil
	}
	if len(fracPart) <= int(decimals) {
		if fracPart != "" {
			fracValue, err := parseDigits(fracPart)
			if err != nil {
				return 0, err
			}
			pad := int(decimals) - len(fracPart)
			if pad > 0 {
				fracValue, err = safeMul(fracValue, pow10(pad))
				if err != nil {
					return 0, err
				}
			}
			amount, err = safeAdd(amount, fracValue)
			if err != nil {
				return 0, err
			}
		}
		return amount * sign, nil
	}
	decimalsInt := int(decimals)
	mainFrac := fracPart[:decimalsInt]
	remainder := fracPart[decimalsInt:]
	fracValue, err := parseDigits(mainFrac)
	if err != nil {
		return 0, err
	}
	amount, err = safeAdd(amount, fracValue)
	if err != nil {
		return 0, err
	}
	if shouldRound(remainder, rounding, amount) {
		amount, err = safeAdd(amount, 1)
		if err != nil {
			return 0, err
		}
	}
	return amount * sign, nil
}

func parseDigits(value string) (int64, error) {
	if value == "" {
		return 0, nil
	}
	var result int64
	for i := 0; i < len(value); i++ {
		ch := value[i]
		if ch < '0' || ch > '9' {
			return 0, ErrInvalidAmount
		}
		digit := int64(ch - '0')
		if result > (math.MaxInt64-digit)/10 {
			return 0, ErrAmountOverflow
		}
		result = result*10 + digit
	}
	return result, nil
}

func validateDigits(value string) error {
	for i := 0; i < len(value); i++ {
		ch := value[i]
		if ch < '0' || ch > '9' {
			return ErrInvalidAmount
		}
	}
	return nil
}

func shouldRound(remainder string, mode RoundingMode, truncated int64) bool {
	if remainder == "" || mode == RoundTruncate {
		return false
	}
	first := remainder[0]
	if first > '5' {
		return true
	}
	if first < '5' {
		return false
	}
	for i := 1; i < len(remainder); i++ {
		if remainder[i] != '0' {
			return true
		}
	}
	if mode == RoundHalfUp {
		return true
	}
	if mode == RoundBankers {
		return truncated%2 != 0
	}
	return false
}

func safeMul(a, b int64) (int64, error) {
	if a == 0 || b == 0 {
		return 0, nil
	}
	if a > math.MaxInt64/b {
		return 0, ErrAmountOverflow
	}
	return a * b, nil
}

func safeAdd(a, b int64) (int64, error) {
	if b > math.MaxInt64-a {
		return 0, ErrAmountOverflow
	}
	return a + b, nil
}
