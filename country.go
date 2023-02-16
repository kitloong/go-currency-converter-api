package currconv

import "net/url"

// Country is the result of Country API.
type Country struct {
	Results map[string]struct {
		ID             string `json:"id"`
		Alpha3         string `json:"alpha3"`
		CurrencyID     string `json:"currencyId"`
		CurrencyName   string `json:"currencyName"`
		CurrencySymbol string `json:"currencySymbol"`
		Name           string `json:"name"`
	} `json:"results"`
}

// Countries returns a list of countries.
func (a *API) Countries() (result *Country, err error) {
	return call[Country](a, true, "countries", func(q url.Values) error { return nil })
}
