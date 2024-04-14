package currency

import (
	"fmt"
	"math"

	"github.com/bqdanh/money_transfer/internal/entities/exceptions"
)

const VND Unit = "VND"

func init() {
	RegisterCurrencyParser(VND, VNDCurrencyParser)
}

func IsFractionalZero(value float64) bool {
	// Check if the fractional part is close to zero within a small epsilon
	return math.Abs(math.Mod(value, 1)) < 1e-10
}

func VNDCurrencyParser(amount float64) (Amount, error) {
	if amount < 0 {
		err := exceptions.NewInvalidArgumentError("amount", "amount must be positive", map[string]interface{}{
			"amount": amount,
		})
		return Amount{}, err
	}
	if !IsFractionalZero(amount) {
		return Amount{}, fmt.Errorf("invalid fractional part: %f", amount)
	}
	return Amount{
		Currency:        VND,
		Amount:          int64(amount),
		FractionalUnits: 0,
	}, nil
}
