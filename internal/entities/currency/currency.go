package currency

import (
	"fmt"
	"sync"
)

type Unit string

type Amount struct {
	Currency        Unit  `json:"currency"`
	Amount          int64 `json:"amount"`
	FractionalUnits int64 `json:"fractional_units"`
}

func (a Amount) Float64() float64 {
	return float64(a.Amount) + float64(a.FractionalUnits)/100
}

func (a Amount) String() string {
	return fmt.Sprintf("%.2f", a.Float64())
}

type AmountParser func(amount float64) (Amount, error)

var (
	currencyStringMapping = map[string]Unit{}
	currencyParserMapping = map[Unit]AmountParser{}
	rwLock                = sync.RWMutex{}
)

func RegisterCurrencyParser(currency Unit, parser AmountParser) {
	rwLock.Lock()
	defer rwLock.Unlock()
	if _, ok := currencyParserMapping[currency]; ok {
		panic(fmt.Errorf("currency(%s) already registered", currency))
	}
	currencyParserMapping[currency] = parser
	currencyStringMapping[string(currency)] = currency
}

func GetCurrencyUnit(currency string) (Unit, error) {
	rwLock.RLock()
	defer rwLock.RUnlock()
	if unit, ok := currencyStringMapping[currency]; ok {
		return unit, nil
	}
	return "", fmt.Errorf("currency(%s) not found", currency)
}

func GetCurrencyParser(currency Unit) (AmountParser, error) {
	rwLock.RLock()
	defer rwLock.RUnlock()
	if parser, ok := currencyParserMapping[currency]; ok {
		return parser, nil
	}
	return nil, fmt.Errorf("currency(%s) not found", currency)
}

func FromFloat64(amount float64, currency Unit) (Amount, error) {
	parser, err := GetCurrencyParser(currency)
	if err != nil {
		return Amount{}, fmt.Errorf("get currency(%s) parser: %w", currency, err)
	}
	v, err := parser(amount)
	if err != nil {
		return Amount{}, fmt.Errorf("parse currency(%s) amount: %w", currency, err)
	}
	return v, nil
}

func (a Amount) IsGt(v float64) (bool, error) {
	parser, err := GetCurrencyParser(a.Currency)
	if err != nil {
		return false, fmt.Errorf("get currency(%s) parser: %w", a.Currency, err)
	}
	amount, err := parser(v)
	if err != nil {
		return false, fmt.Errorf("parse currency(%s) amount: %w", a.Currency, err)
	}
	if a.Amount > amount.Amount {
		return true, nil
	}
	if a.Amount == amount.Amount && a.FractionalUnits > amount.FractionalUnits {
		return true, nil
	}
	return false, nil
}

func (a Amount) IsLt(v float64) (bool, error) {
	parser, err := GetCurrencyParser(a.Currency)
	if err != nil {
		return false, fmt.Errorf("get currency(%s) parser: %w", a.Currency, err)
	}
	amount, err := parser(v)
	if err != nil {
		return false, fmt.Errorf("parse currency(%s) amount: %w", a.Currency, err)
	}
	if a.Amount < amount.Amount {
		return true, nil
	}
	if a.Amount == amount.Amount && a.FractionalUnits < amount.FractionalUnits {
		return true, nil
	}
	return false, nil
}
