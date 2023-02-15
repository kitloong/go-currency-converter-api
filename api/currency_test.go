package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPI_Currencies(t *testing.T) {
	tests := []struct {
		name     string
		respJSON []byte
		error    error
	}{
		{
			"Success response",
			[]byte(`{
				"results": {
					"MYR": {
						"currencyName": "Malaysian Ringgit",
						"currencySymbol": "RM",
						"id": "MYR"
					},
					"USD": {
						"currencyName": "United States Dollar",
						"currencySymbol": "$",
						"id": "USD"
					}
				}
			}`),
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, err := w.Write(tt.respJSON)
				if err != nil {
					return
				}

				assert.Equal(t, "/api/v1/currencies", r.URL.Path)

				q := url.Values{}
				q.Add("apiKey", "key")

				assert.Equal(t, q, r.URL.Query())
			}))
			defer ts.Close()

			api := NewAPI(Config{
				BaseURL: ts.URL,
				APIKey:  "key",
				Version: "v1",
			})

			convert, err := api.Currencies()
			if err != nil {
				assert.Equal(t, tt.error, err)
				return
			}

			expected := &Currency{}

			err = json.Unmarshal(tt.respJSON, expected)
			if err != nil {
				return
			}

			assert.Equal(t, expected, convert)
		})
	}
}
