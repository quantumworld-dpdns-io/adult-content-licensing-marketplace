package payment

import (
	"errors"
	"strings"
)

var ErrNonCryptoCurrency = errors.New("only crypto-denominated currencies are allowed")

var allowed = map[string]struct{}{
	"USDC": {},
	"USDT": {},
	"ETH":  {},
}

func EnsureCryptoCurrency(currency string) error {
	c := strings.ToUpper(strings.TrimSpace(currency))
	if _, ok := allowed[c]; !ok {
		return ErrNonCryptoCurrency
	}
	return nil
}
