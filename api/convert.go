package api

import (
	"errors"
	"net/url"
	"strings"
	"time"
)

// ConvertRequest contains request fields of Convert and ConvertCompact API.
type ConvertRequest struct {
	// Q is the currency conversion parameter in "[FROM]_[TO]" format.
	// To get MYR -> USD conversion rate, uses "MYR_USD".
	// A multiple Q request will return multiple conversion in a single request.
	Q []string
}

// Convert is the result of the Convert API.
type Convert struct {
	Query struct {
		Count int `json:"count"`
	} `json:"query"`
	Results map[string]struct {
		ID  string  `json:"id"`
		Val float32 `json:"val"`
		To  string  `json:"to"`
		Fr  string  `json:"fr"`
	} `json:"results"`
}

// ConvertCompact is the compact result of the ConvertCompact API.
type ConvertCompact map[string]float32

// ConvertHistoricalRequest contains request fields of ConvertHistorical and ConvertHistoricalCompact API.
type ConvertHistoricalRequest struct {
	// Q is the same with ConvertRequest's Q.
	Q []string
	// Date is the target historical date that you want to request from.
	Date time.Time
	// EndDate form a date range with Date, returns historical data in within the date range.
	EndDate time.Time
}

// ConvertHistorical is the result of ConvertHistorical API.
type ConvertHistorical struct {
	Query struct {
		Count int `json:"count"`
	} `json:"query"`
	Date    string `json:"date"`
	EndDate string `json:"endDate,omitempty"`
	Results map[string]struct {
		ID  string             `json:"id"`
		To  string             `json:"to"`
		Fr  string             `json:"fr"`
		Val map[string]float32 `json:"val"`
	} `json:"results"`
}

// ConvertHistoricalCompact is the compact result of ConvertHistoricalCompact API.
type ConvertHistoricalCompact map[string]map[string]float32

// Convert returns the currency conversion rate with `[FROM]_[TO]` request.
func (a *API) Convert(req ConvertRequest) (result *Convert, err error) {
	return call[Convert](a, true, "convert", func(q url.Values) error {
		if len(req.Q) == 0 {
			return errors.New("`Q` require at least one currency conversion")
		}

		q.Add("q", strings.Join(req.Q, ","))
		return nil
	})
}

// ConvertCompact returns conversion result with compact mode.
func (a *API) ConvertCompact(req ConvertRequest) (result ConvertCompact, err error) {
	r, err := call[ConvertCompact](a, true, "convert", func(q url.Values) error {
		if len(req.Q) == 0 {
			return errors.New("`Q` require at least one currency conversion")
		}

		q.Add("compact", "ultra")
		q.Add("q", strings.Join(req.Q, ","))
		return nil
	})
	if err != nil {
		return ConvertCompact{}, err
	}

	return *r, nil
}

// ConvertHistorical returns historical currency conversion rate data with target date or date range.
func (a *API) ConvertHistorical(req ConvertHistoricalRequest) (result *ConvertHistorical, err error) {
	return call[ConvertHistorical](a, true, "convert", func(q url.Values) error {
		if len(req.Q) == 0 {
			return errors.New("`Q` require at least one currency conversion")
		}

		if req.Date.IsZero() {
			return errors.New("`Date` is required")
		}

		q.Add("q", strings.Join(req.Q, ","))
		q.Add("date", req.Date.Format("2006-01-02"))
		if !req.EndDate.IsZero() {
			q.Add("endDate", req.EndDate.Format("2006-01-02"))
		}
		return nil
	})
}

// ConvertHistoricalCompact returns historical data with compact mode.
func (a *API) ConvertHistoricalCompact(req ConvertHistoricalRequest) (result ConvertHistoricalCompact, err error) {
	r, err := call[ConvertHistoricalCompact](a, true, "convert", func(q url.Values) error {
		if len(req.Q) == 0 {
			return errors.New("`Q` require at least one currency conversion")
		}

		if req.Date.IsZero() {
			return errors.New("`Date` is required")
		}

		q.Add("compact", "ultra")
		q.Add("q", strings.Join(req.Q, ","))
		q.Add("date", req.Date.Format("2006-01-02"))
		if !req.EndDate.IsZero() {
			q.Add("endDate", req.EndDate.Format("2006-01-02"))
		}
		return nil
	})
	if err != nil {
		return ConvertHistoricalCompact{}, err
	}

	return *r, nil
}
