package currconv

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAPI_Convert(t *testing.T) {
	tests := []struct {
		name     string
		req      ConvertRequest
		respJSON []byte
		error    error
	}{
		{
			"Success",
			ConvertRequest{
				Q: []string{"USD_MYR"},
			},
			[]byte(`{
				"query": {
					"count": 1
				},
				"results": {
					"USD_MYR": {
						"id": "USD_MYR",
						"val": 4.348493,
						"to": "MYR",
						"fr": "USD"
					}
				}
			}`),
			nil,
		},
		{
			"Success with multiple values",
			ConvertRequest{
				Q: []string{"USD_MYR", "MYR_USD"},
			},
			[]byte(`{
				"query": {
					"count": 2
				},
				"results": {
					"USD_MYR": {
						"id": "USD_MYR",
						"val": 4.348493,
						"to": "MYR",
						"fr": "USD"
					},
					"MYR_USD": {
						"id": "MYR_USD",
						"val": 0.229964,
						"to": "USD",
						"fr": "MYR"
					}
				}
			}`),
			nil,
		},
		{
			"`Q` is required",
			ConvertRequest{},
			[]byte(``),
			errors.New("`Q` require at least one currency conversion"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, err := w.Write(tt.respJSON)
				if err != nil {
					return
				}

				assert.Equal(t, "/api/v1/convert", r.URL.Path)

				q := url.Values{}
				q.Add("apiKey", "key")
				q.Add("q", strings.Join(tt.req.Q, ","))

				assert.Equal(t, q, r.URL.Query())
			}))
			defer ts.Close()

			api := NewAPI(Config{
				BaseURL: ts.URL,
				APIKey:  "key",
				Version: "v1",
			})

			convert, err := api.Convert(tt.req)
			if err != nil {
				assert.Equal(t, tt.error, err)
				return
			}

			expected := &Convert{}

			err = json.Unmarshal(tt.respJSON, expected)
			if err != nil {
				return
			}

			assert.Equal(t, expected, convert)
		})
	}
}

func TestAPI_ConvertCompact(t *testing.T) {
	tests := []struct {
		name     string
		req      ConvertRequest
		respJSON []byte
		error    error
	}{
		{
			"Success",
			ConvertRequest{
				Q: []string{"USD_MYR"},
			},
			[]byte(`{
 				"USD_MYR": 4.348493
			}`),
			nil,
		},
		{
			"Success with multiple values",
			ConvertRequest{
				Q: []string{"USD_MYR", "MYR_USD"},
			},
			[]byte(`{
 				"USD_MYR": 4.348493,
				"MYR_USD": 0.229964
			}`),
			nil,
		},
		{
			"`Q` is required",
			ConvertRequest{},
			[]byte(``),
			errors.New("`Q` require at least one currency conversion"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, err := w.Write(tt.respJSON)
				if err != nil {
					return
				}

				assert.Equal(t, "/api/v1/convert", r.URL.Path)

				q := url.Values{}
				q.Add("apiKey", "key")
				q.Add("compact", "ultra")
				q.Add("q", strings.Join(tt.req.Q, ","))

				assert.Equal(t, q, r.URL.Query())
			}))
			defer ts.Close()

			api := NewAPI(Config{
				BaseURL: ts.URL,
				APIKey:  "key",
				Version: "v1",
			})

			convert, err := api.ConvertCompact(tt.req)
			if err != nil {
				assert.Equal(t, tt.error, err)
				return
			}

			expected := ConvertCompact{}

			err = json.Unmarshal(tt.respJSON, &expected)
			if err != nil {
				return
			}

			assert.Equal(t, expected, convert)
		})
	}
}

func TestAPI_ConvertHistorical(t *testing.T) {
	tests := []struct {
		name     string
		req      ConvertHistoricalRequest
		respJSON []byte
		error    error
	}{
		{
			"Success",
			ConvertHistoricalRequest{
				Q:    []string{"USD_MYR"},
				Date: time.Date(2023, 2, 14, 0, 0, 0, 0, time.UTC),
			},
			[]byte(`{
				"query": {
					"count": 1
				},
				"date": "2023-02-14",
				"results": {
					"USD_MYR": {
						"id": "USD_MYR",
						"fr": "USD",
						"to": "MYR",
						"val": {
							"2023-02-14": 4.348493
						}
					}
				}
			}`),
			nil,
		},
		{
			"Success with multiple values",
			ConvertHistoricalRequest{
				Q:    []string{"USD_MYR", "MYR_USD"},
				Date: time.Date(2023, 2, 14, 0, 0, 0, 0, time.UTC),
			},
			[]byte(`{
				"query": {
					"count": 2
				},
				"date": "2023-02-14",
				"results": {
					"USD_MYR": {
						"id": "USD_MYR",
						"fr": "USD",
						"to": "MYR",
						"val": {
							"2023-02-14": 4.348493
						}
					},
					"MYR_USD": {
						"id": "MYR_USD",
						"fr": "MYR",
						"to": "USD",
						"val": {
							"2023-02-14": 0.229964
						}
					}
				}
			}`),
			nil,
		},
		{
			"Success date range",
			ConvertHistoricalRequest{
				Q:       []string{"USD_MYR"},
				Date:    time.Date(2023, 2, 14, 0, 0, 0, 0, time.UTC),
				EndDate: time.Date(2023, 2, 16, 0, 0, 0, 0, time.UTC),
			},
			[]byte(`{
				"query": {
					"count": 1
				},
				"date": "2023-02-14",
				"endDate": "2023-02-16",
				"results": {
					"USD_MYR": {
						"id": "USD_MYR",
						"fr": "USD",
						"to": "MYR",
						"val": {
							"2023-02-14": 4.369895,
							"2023-02-15": 4.379895,
							"2023-02-16": 4.359895
						}
					}
				}
			}`),
			nil,
		},
		{
			"Success date range with multiple values",
			ConvertHistoricalRequest{
				Q:       []string{"USD_MYR", "MYR_USD"},
				Date:    time.Date(2023, 2, 14, 0, 0, 0, 0, time.UTC),
				EndDate: time.Date(2023, 2, 16, 0, 0, 0, 0, time.UTC),
			},
			[]byte(`{
				"query": {
					"count": 2
				},
				"date": "2023-02-14",
				"endDate": "2023-02-16",
				"results": {
					"USD_MYR": {
						"id": "USD_MYR",
						"fr": "USD",
						"to": "MYR",
						"val": {
							"2023-02-14": 4.369895,
							"2023-02-15": 4.379895,
							"2023-02-16": 4.359895
						}
					},
					"MYR_USD": {
						"id": "MYR_USD",
						"fr": "MYR",
						"to": "USD",
						"val": {
							"2023-02-14": 0.227824,
							"2023-02-15": 0.227914,
							"2023-02-16": 0.227927
						}
					}
				}
			}`),
			nil,
		},
		{
			"`Q` is required",
			ConvertHistoricalRequest{},
			[]byte(``),
			errors.New("`Q` require at least one currency conversion"),
		},
		{
			"`Date` is required",
			ConvertHistoricalRequest{
				Q: []string{"USD_MYR", "MYR_USD"},
			},
			[]byte(``),
			errors.New("`Date` is required"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, err := w.Write(tt.respJSON)
				if err != nil {
					return
				}

				assert.Equal(t, "/api/v1/convert", r.URL.Path)

				q := url.Values{}
				q.Add("apiKey", "key")
				q.Add("q", strings.Join(tt.req.Q, ","))
				q.Add("date", tt.req.Date.Format("2006-01-02"))
				if !tt.req.EndDate.IsZero() {
					q.Add("endDate", tt.req.EndDate.Format("2006-01-02"))
				}

				assert.Equal(t, q, r.URL.Query())
			}))
			defer ts.Close()

			api := NewAPI(Config{
				BaseURL: ts.URL,
				APIKey:  "key",
				Version: "v1",
			})

			convert, err := api.ConvertHistorical(tt.req)
			if err != nil {
				assert.Equal(t, tt.error, err)
				return
			}

			expected := &ConvertHistorical{}

			err = json.Unmarshal(tt.respJSON, expected)
			if err != nil {
				return
			}

			assert.Equal(t, expected, convert)
		})
	}
}

func TestAPI_ConvertHistoricalCompact(t *testing.T) {
	tests := []struct {
		name     string
		req      ConvertHistoricalRequest
		respJSON []byte
		error    error
	}{
		{
			"Success",
			ConvertHistoricalRequest{
				Q:    []string{"USD_MYR"},
				Date: time.Date(2023, 2, 14, 0, 0, 0, 0, time.UTC),
			},
			[]byte(`{
				"USD_MYR": {
					"2023-02-14": 4.348493
				}
			}`),
			nil,
		},
		{
			"Success with multiple values",
			ConvertHistoricalRequest{
				Q:    []string{"USD_MYR", "MYR_USD"},
				Date: time.Date(2023, 2, 14, 0, 0, 0, 0, time.UTC),
			},
			[]byte(`{
				"USD_MYR": {
					"2023-02-14": 4.348493
				},
				"MYR_USD": {
					"2023-02-14": 0.229964
				}
			}`),
			nil,
		},
		{
			"Success date range",
			ConvertHistoricalRequest{
				Q:       []string{"USD_MYR"},
				Date:    time.Date(2023, 2, 14, 0, 0, 0, 0, time.UTC),
				EndDate: time.Date(2023, 2, 16, 0, 0, 0, 0, time.UTC),
			},
			[]byte(`{
				"USD_MYR": {
					"2023-02-14": 4.369895,
					"2023-02-15": 4.379895,
					"2023-02-16": 4.359895
				}
			}`),
			nil,
		},
		{
			"Success date range with multiple values",
			ConvertHistoricalRequest{
				Q:       []string{"USD_MYR", "MYR_USD"},
				Date:    time.Date(2023, 2, 14, 0, 0, 0, 0, time.UTC),
				EndDate: time.Date(2023, 2, 16, 0, 0, 0, 0, time.UTC),
			},
			[]byte(`{
				"USD_MYR": {
					"2023-02-14": 4.369895,
					"2023-02-15": 4.379895,
					"2023-02-16": 4.359895
				},
				"MYR_USD": {
					"2023-02-14": 0.217824,
					"2023-02-15": 0.217914,
					"2023-02-16": 0.217927
				}
			}`),
			nil,
		},
		{
			"`Q` is required",
			ConvertHistoricalRequest{},
			[]byte(``),
			errors.New("`Q` require at least one currency conversion"),
		},
		{
			"`Date` is required",
			ConvertHistoricalRequest{
				Q: []string{"USD_MYR", "MYR_USD"},
			},
			[]byte(``),
			errors.New("`Date` is required"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, err := w.Write(tt.respJSON)
				if err != nil {
					return
				}

				assert.Equal(t, "/api/v1/convert", r.URL.Path)

				q := url.Values{}
				q.Add("apiKey", "key")
				q.Add("compact", "ultra")
				q.Add("q", strings.Join(tt.req.Q, ","))
				q.Add("date", tt.req.Date.Format("2006-01-02"))
				if !tt.req.EndDate.IsZero() {
					q.Add("endDate", tt.req.EndDate.Format("2006-01-02"))
				}

				assert.Equal(t, q, r.URL.Query())
			}))
			defer ts.Close()

			api := NewAPI(Config{
				BaseURL: ts.URL,
				APIKey:  "key",
				Version: "v1",
			})

			convert, err := api.ConvertHistoricalCompact(tt.req)
			if err != nil {
				assert.Equal(t, tt.error, err)
				return
			}

			expected := ConvertHistoricalCompact{}

			err = json.Unmarshal(tt.respJSON, &expected)
			if err != nil {
				return
			}

			assert.Equal(t, expected, convert)
		})
	}
}
