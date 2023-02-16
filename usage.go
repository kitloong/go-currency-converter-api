package currconv

import (
	"net/url"
	"time"
)

// Usage is the result of Usage API.
type Usage struct {
	Timestamp time.Time `json:"timestamp"`
	Usage     int       `json:"usage"`
}

// Usage returns your current API usage.
func (a *API) Usage() (result *Usage, err error) {
	return call[Usage](a, false, "others/usage", func(q url.Values) error { return nil })
}
