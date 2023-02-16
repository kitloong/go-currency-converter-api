package currconv

import (
	"net/url"
)

// Currency is the result of Currencies API.
type Currency struct {
	Results map[string]struct {
		ID             string `json:"id"`
		CurrencyName   string `json:"currencyName"`
		CurrencySymbol string `json:"currencySymbol"`
	} `json:"results"`
}

// Currencies returns a list of currencies.
func (a *API) Currencies() (result *Currency, err error) {
	return call[Currency](a, true, "currencies", func(q url.Values) error { return nil })
}
