package currconv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAPI(t *testing.T) {
	api := NewAPI(Config{
		BaseURL: "baseURL",
		APIKey:  "APIKey",
		Version: "Version",
	})

	assert.Equal(t, api.config.BaseURL, "baseURL")
	assert.Equal(t, api.config.APIKey, "APIKey")
	assert.Equal(t, api.config.Version, "Version")
}

func TestAPIError(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		errorMsg string
	}{
		{
			"URL parse error",
			"http://[::1]a",
			"parse \"http://[::1]a\": invalid port \"a\" after host",
		},
		{
			"HTTP Get error",
			"/error/",
			"Get \"/error/api/Version/convert?apiKey=APIKey&q=MYR_USD\": unsupported protocol scheme \"\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := NewAPI(Config{
				BaseURL: tt.baseURL,
				APIKey:  "APIKey",
				Version: "Version",
			})
			_, err := api.Convert(ConvertRequest{Q: []string{"MYR_USD"}})
			if err != nil {
				assert.EqualError(t, err, tt.errorMsg)
			}
		})
	}
}
