package currconv

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPI_Countries(t *testing.T) {
	tests := []struct {
		name     string
		respJSON []byte
		error    error
	}{
		{
			"Success response",
			[]byte(`{
				"results": {
					"MY": {
						"alpha3": "MYS",
						"currencyId": "MYR",
						"currencyName": "Malaysian ringgit",
						"currencySymbol": "RM",
						"id": "MY",
						"name": "Malaysia"
					},
					"US": {
						"alpha3": "USA",
						"currencyId": "USD",
						"currencyName": "United States dollar",
						"currencySymbol": "$",
						"id": "US",
						"name": "United States of America"
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

				assert.Equal(t, "/api/v1/countries", r.URL.Path)

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

			convert, err := api.Countries()
			if err != nil {
				assert.Equal(t, tt.error, err)
				return
			}

			expected := &Country{}

			err = json.Unmarshal(tt.respJSON, expected)
			if err != nil {
				return
			}

			assert.Equal(t, expected, convert)
		})
	}
}
