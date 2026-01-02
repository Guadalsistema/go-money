package currency_test

import (
	"testing"

	money "github.com/Guadalsistema/go-money"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	a := money.USD(1000) // $10.00
	b := money.USD(2500) // $25.00

	result, err := money.Add(a, b)
	assert.NoError(t, err)
	assert.Equal(t, money.USD(3500), result) // $35.00
}

func TestSubtract(t *testing.T) {
	a := money.EUR(5000) // €50.00
	b := money.EUR(2000) // €20.00

	result, err := money.Subtract(a, b)
	assert.NoError(t, err)
	assert.Equal(t, money.EUR(3000), result) // €30.00
}

func TestAddDifferentCurrencies(t *testing.T) {
	a := money.USD(1000) // $10.00
	b := money.EUR(2500) // €25.00

	_, err := money.Add(a, b)
	assert.ErrorIs(t, err, money.ErrDifferentCurrency)
}
