package currency_test

import (
	"encoding/json"
	"testing"

	money "github.com/Guadalsistema/go-money"
	"github.com/stretchr/testify/assert"
)

func TestCurrencyJSONRoundTrip(t *testing.T) {
	original := money.USD(1234)

	payload, err := json.Marshal(original)
	assert.NoError(t, err)

	var decoded money.Currency
	err = json.Unmarshal(payload, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, original, decoded)
	assert.EqualValues(t, money.USDDecimals, decoded.Decimals())
}

func TestCurrencyJSONInferDecimals(t *testing.T) {
	payload := []byte(`{"amount":500,"currency":"JPY"}`)

	var decoded money.Currency
	err := json.Unmarshal(payload, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, int64(500), decoded.Amount)
	assert.Equal(t, "JPY", decoded.Currency)
	assert.EqualValues(t, money.JPYDecimals, decoded.Decimals())
}

func TestCurrencyJSONUnsupportedCurrency(t *testing.T) {
	payload := []byte(`{"amount":500,"currency":"ZZZ"}`)

	var decoded money.Currency
	err := json.Unmarshal(payload, &decoded)
	assert.ErrorIs(t, err, money.ErrUnsupportedCurrency)
}
