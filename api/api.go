package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

// Config of the API.
type Config struct {
	BaseURL string
	Version string
	APIKey  string
}

// API is the wrapper implementation of CurrencyConverterAPI.
type API struct {
	config Config
}

type Error struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}

// NewAPI create and return an API.
func NewAPI(config Config) *API {
	return &API{
		config,
	}
}

type response interface {
	Convert | ConvertCompact | ConvertHistorical | ConvertHistoricalCompact | Currency | Country | Usage
}

// call is a function used by all APIs to request CurrencyConverterAPI.
// This function will execute `handler` which consist unique logic from the caller.
func call[T response](a *API, shouldPrefixAPIPath bool, path string, handler func(q url.Values) error) (result *T, err error) {
	u, err := url.Parse(a.config.BaseURL)
	if err != nil {
		return nil, err
	}

	if shouldPrefixAPIPath {
		u = u.JoinPath("api").JoinPath(a.config.Version)
	}

	u = u.JoinPath(path)

	query := u.Query()
	query.Add("apiKey", a.config.APIKey)

	err = handler(query)
	if err != nil {
		return nil, err
	}

	u.RawQuery = query.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, parseError(resp)
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return
}

// parseError uses `json.Unmarshal` to returns Error whenever it is possible.
func parseError(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		return err
	}

	e := Error{}
	err = json.Unmarshal(body, &e)
	if err != nil {
		return errors.New(string(body))
	}

	return errors.New(e.Error)
}
