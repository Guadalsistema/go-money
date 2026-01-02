package currency

func Add(a, b Currency) (Currency, error) {
	if a.Currency != b.Currency {
		return Currency{}, ErrDifferentCurrency
	}
	return Currency{
		Amount:   a.Amount + b.Amount,
		Currency: a.Currency,
		decimals: a.decimals,
	}, nil
}

func Subtract(a, b Currency) (Currency, error) {
	if a.Currency != b.Currency {
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
