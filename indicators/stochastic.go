package indicators

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/DefinitelyNotAGoat/twelvedata/model"
	"github.com/pkg/errors"
)

// StochasticIndicator - the Indicator value for IndicatorMeta specific for Stochastic
type StochasticIndicator struct {
	Name        string `json:"name"`
	FastKPeriod int    `json:"fast_k_period"`
	SlowKPeriod int    `json:"slow_k_period"`
	SlowDPeriod int    `json:"slow_d_period"`
	SlowKmaType string `json:"slow_kma_type"`
	SlowDmaType string `json:"slow_dma_type"`
}

// StochasticValue - the Indicator value for IndicatorResponse specific for Stochastic
type StochasticValue struct {
	Datetime time.Time `json:"datetime"`
	SlowK    float64   `json:"slow_k"`
	SlowD    float64   `json:"slow_d"`
}

// UnmarshalJSON - unmarshal's StochasticValue to a more consumable type
func (s *StochasticValue) UnmarshalJSON(v []byte) error {
	type Value struct {
		Datetime string `json:"datetime"`
		SlowK    string `json:"slow_k"`
		SlowD    string `json:"slow_d"`
	}

	var value Value
	if err := json.Unmarshal(v, &value); err != nil {
		return err
	}

	dateTime, err := time.Parse(model.GetTimeFormatFromString(value.Datetime), value.Datetime)
	if err != nil {
		return errors.Wrap(err, "failed to parse value date time into go time")
	}
	s.Datetime = dateTime

	slowK, err := strconv.ParseFloat(value.SlowK, 64)
	if err != nil {
		return errors.Wrap(err, "failed to parse value slowk into float")
	}

	slowD, err := strconv.ParseFloat(value.SlowD, 64)
	if err != nil {
		return errors.Wrap(err, "failed to parse value slowD into float")
	}

	s.SlowK = slowK
	s.SlowD = slowD

	return nil
}

// StochasticOptions - options for calling the twelvedata stoch endpoint: https://twelvedata.com/docs#stoch
type StochasticOptions struct {
	IndicatorOptions
	FastKPeriod int
	SlowDPeriod int
	SlowDMAType string
	SlowKPeriod int
	SlowKMAType string
}

func (s StochasticOptions) params(u *url.URL, urlValues url.Values) {
	urlValues = s.IndicatorOptions.params(u, urlValues)

	if s.FastKPeriod > 0 {
		urlValues.Add("fast_k_period", strconv.Itoa(s.FastKPeriod))
	}

	if s.SlowDPeriod > 0 {
		urlValues.Add("slow_d_period", strconv.Itoa(s.SlowDPeriod))
	}

	if s.SlowDMAType == "" {
		urlValues.Add("slow_dma_type", s.SlowDMAType)
	}

	if s.SlowKPeriod > 0 {
		urlValues.Add("slow_k_period", strconv.Itoa(s.SlowKPeriod))
	}

	if s.SlowKMAType == "" {
		urlValues.Add("slow_kma_type", s.SlowKMAType)
	}

	u.RawQuery = urlValues.Encode()
}

func stochastic[Response IndicatorResponse[StochasticValue, StochasticIndicator]](symbol string, interval model.Interval, apiKey string, getFn getFn, opts StochasticOptions) (Response, error) {
	u, err := url.Parse(fmt.Sprintf("%s/stoch", baseURI))
	if err != nil {
		return Response{}, errors.Wrapf(err, "failed to parse base URL '%s'", baseURI)
	}

	opts.params(u, url.Values{
		"symbol":   {symbol},
		"interval": {string(interval)},
		"apikey":   {apiKey},
	})

	body, err := getFn(u)
	if err != nil {
		return Response{}, err
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return Response{}, err
	}

	return response, nil
}
