package currency_test

import (
	"testing"

	money "github.com/Guadalsistema/go-money"
	"github.com/stretchr/testify/assert"
)

func TestParseAmountBasic(t *testing.T) {
	got, err := money.ParseAmount("EUR", "12.34")
	assert.NoError(t, err)
	assert.Equal(t, int64(1234), got.Amount)
	assert.Equal(t, "EUR", got.Currency)
	assert.EqualValues(t, money.EURDecimals, got.Decimals())
}

func TestParseAmountNegative(t *testing.T) {
	got, err := money.ParseAmount("USD", "-1.23")
	assert.NoError(t, err)
	assert.Equal(t, int64(-123), got.Amount)
}

func TestParseAmountRoundingHalfUp(t *testing.T) {
	got, err := money.ParseAmount("USD", "1.235", money.WithRounding(money.RoundHalfUp))
	assert.NoError(t, err)
	assert.Equal(t, int64(124), got.Amount)
}

func TestParseAmountRoundingBankers(t *testing.T) {
	odd, err := money.ParseAmount("USD", "1.235", money.WithRounding(money.RoundBankers))
	assert.NoError(t, err)
	assert.Equal(t, int64(124), odd.Amount)

	even, err := money.ParseAmount("USD", "1.245", money.WithRounding(money.RoundBankers))
	assert.NoError(t, err)
	assert.Equal(t, int64(124), even.Amount)
}

func TestParseAmountOrZero(t *testing.T) {
	got := money.ParseAmountOrZero("EUR", "oops")
	assert.Equal(t, int64(0), got.Amount)
	assert.Equal(t, "EUR", got.Currency)
	assert.EqualValues(t, money.EURDecimals, got.Decimals())
}

func TestDecimalString(t *testing.T) {
	amount := money.USD(-1234)
	assert.Equal(t, "-12.34", amount.DecimalString())
}
