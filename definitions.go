package currency

// EUR creates a Currency instance for Euro currency.
func EUR(amount int64) Currency {
	return Currency{
		Amount:   amount,
		Currency: "EUR",
		decimals: EURDecimals,
	}
}

// USD creates a Currency instance for US Dollar currency.
func USD(amount int64) Currency {
	return Currency{
		Amount:   amount,
		Currency: "USD",
		decimals: USDDecimals,
	}
}

// Yuan creates a Currency instance for Chinese Yuan currency.
func Yuan(amount int64) Currency {
	return Currency{
		Amount:   amount,
		Currency: "CNY",
		decimals: 2,
	}
}

// MorroccanDirham creates a Currency instance for Moroccan Dirham currency.
func MorroccanDirham(amount int64) Currency {
	return Currency{
		Amount:   amount,
		Currency: "MAD",
		decimals: 2,
	}
}
