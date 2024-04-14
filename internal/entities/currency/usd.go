package currency

import (
	"github.com/bqdanh/money_transfer/internal/entities/exceptions"
)

const USD Unit = "USD"

func init() {
	RegisterCurrencyParser(USD, USDCurrencyParser)
}

func USDCurrencyParser(amount float64) (Amount, error) {
	if amount < 0 {
		err := exceptions.NewInvalidArgumentError("amount", "amount must be positive", map[string]interface{}{
			"amount": amount,
		})
		return Amount{}, err
	}
	usdAmount := float64(amount) * 100
	if !IsFractionalZero(usdAmount) {
		return Amount{}, exceptions.NewInvalidArgumentError("amount", "invalid fractional part", map[string]interface{}{
			"amount": amount,
		})
	}

	dollars := int64(usdAmount) / 100
	cents := int64(usdAmount) % 100

	return Amount{
		Currency:        USD,
		Amount:          dollars,
		FractionalUnits: cents,
	}, nil
}
