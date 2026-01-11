package currency

func Add(a, b Currency) (Currency, error) {
	if a.Currency != b.Currency || a.decimals != b.decimals {
		return Currency{}, ErrDifferentCurrency
	}
	return Currency{
		Amount:   a.Amount + b.Amount,
		Currency: a.Currency,
		decimals: a.decimals,
	}, nil
}

func Subtract(a, b Currency) (Currency, error) {
	if a.Currency != b.Currency || a.decimals != b.decimals {
		return Currency{}, ErrDifferentCurrency
	}
	return Currency{
		Amount:   a.Amount - b.Amount,
		Currency: a.Currency,
		decimals: a.decimals,
	}, nil
}

func minus(a, b Currency) (Currency, error) {
	return Subtract(a, b)
}

func Divide(a Currency, divisor float64) Currency {
	return Currency{
		Amount:   int64(float64(a.Amount) / divisor),
		Currency: a.Currency,
		decimals: a.decimals,
	}
}

func Multiply(a Currency, factor float64) Currency {
	return Currency{
		Amount:   int64(float64(a.Amount) * factor),
		Currency: a.Currency,
		decimals: a.decimals,
	}
}

// Sum adds a list of Currency values, failing on mismatched currency or decimals.
func Sum(values ...Currency) (Currency, error) {
	if len(values) == 0 {
		return Currency{}, nil
	}
	total := values[0]
	for i := 1; i < len(values); i++ {
		if total.Currency != values[i].Currency || total.decimals != values[i].decimals {
			return Currency{}, ErrDifferentCurrency
		}
		total.Amount += values[i].Amount
	}
	return total, nil
}
