package currency_test

import (
	"testing"

	money "github.com/Guadalsistema/go-money"
	"github.com/stretchr/testify/assert"
)

func TestFromCode(t *testing.T) {
	got, err := money.FromCode("JPY", 500)
	assert.NoError(t, err)
	assert.Equal(t, int64(500), got.Amount)
	assert.EqualValues(t, money.JPYDecimals, got.Decimals())
}

func TestRegisterCurrency(t *testing.T) {
	err := money.RegisterCurrency("BTC", 8)
	assert.NoError(t, err)

	decimals, ok := money.DecimalsFor("BTC")
	assert.True(t, ok)
	assert.EqualValues(t, 8, decimals)
}
