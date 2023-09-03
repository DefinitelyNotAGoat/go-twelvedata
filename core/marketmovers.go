package core

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

type Direction string

const (
	Gainers Direction = "gainers"
	Losers  Direction = "losers"
)

type MarketMoversResponse struct {
	Values []MarketMoversValue `json:"values"`
	Status string              `json:"status"`
}

type MarketMoversValue struct {
	Symbol        string  `json:"symbol"`
	Name          string  `json:"name"`
	Exchange      string  `json:"exchange"`
	MicCode       string  `json:"mic_code"`
	Datetime      string  `json:"datetime"`
	Last          float64 `json:"last"`
	High          float64 `json:"high"`
	Low           float64 `json:"low"`
	Volume        float64 `json:"volume"`
	Change        float64 `json:"change"`
	PercentChange float64 `json:"percent_change"`
}

// MarketMoversOptions - options for calling the twelvedata time series endpoint: https://twelvedata.com/docs#market-movers
type MarketMoversOptions struct {
	Direction  Direction
	OutputSize int
	Country    string
}

func (m MarketMoversOptions) params(u *url.URL, urlValues url.Values) {
	if m.Direction != "" {
		urlValues.Add("direction", string(m.Direction))
	}

	if m.OutputSize > 0 {
		urlValues.Add("outputsize", strconv.Itoa(m.OutputSize))
	}

	if m.Country != "" {
		urlValues.Add("country", m.Country)
	}

	u.RawQuery = urlValues.Encode()
}

// MarketMovers - get the biggest winners or losers depending on opts
func (c *client) MarketMovers(opts MarketMoversOptions) (MarketMoversResponse, error) {
	u, err := url.Parse(fmt.Sprintf("%s/time_series", baseURI))
	if err != nil {
		return MarketMoversResponse{}, errors.Wrapf(err, "failed to parse base URL '%s'", baseURI)
	}

	opts.params(u, url.Values{
		"apikey": {c.apiKey},
	})

	body, err := c.getFn(u)
	if err != nil {
		return MarketMoversResponse{}, err
	}

	var response MarketMoversResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return MarketMoversResponse{}, err
	}

	return response, nil
}
