package currconv

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPI_Usage(t *testing.T) {
	tests := []struct {
		name     string
		respJSON []byte
		error    error
	}{
		{
			"Success response",
			[]byte(`{
				"timestamp": "2023-02-14T10:34:42.039Z",
				"usage": 17
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

				assert.Equal(t, "/others/usage", r.URL.Path)

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

			convert, err := api.Usage()
			if err != nil {
				assert.Equal(t, tt.error, err)
				return
			}

			expected := &Usage{}

			err = json.Unmarshal(tt.respJSON, expected)
			if err != nil {
				return
			}

			assert.Equal(t, expected, convert)
		})
	}
}
