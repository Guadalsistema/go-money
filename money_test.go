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

func TestUnmarshalJSONWithoutDecimals(t *testing.T) {
	// JSON without decimals field, as might be produced by older versions or custom marshaling
	jsonData := `{"amount":525,"currency":"EUR"}`

	var m money.Currency
	err := json.Unmarshal([]byte(jsonData), &m)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// For EUR with 2 decimals, amount 525 should represent 5.25
	expected := 5.25
	if got := m.ToFloat(); got != expected {
		t.Errorf("ToFloat() = %v, expected %v (Decimals = %d)", got, expected, m.Decimals())
	}
}
