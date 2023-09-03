package core

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/DefinitelyNotAGoat/twelvedata/model"
	"github.com/pkg/errors"
)

type TimeSeriesResponse struct {
	Meta   Meta    `json:"meta"`
	Values []Value `json:"values"`
	Status string  `json:"status"`
}

type Meta struct {
	Symbol           string `json:"symbol"`
	Interval         string `json:"interval"`
	Currency         string `json:"currency"`
	ExchangeTimezone string `json:"exchange_timezone"`
	Exchange         string `json:"exchange"`
	MicCode          string `json:"mic_code"`
	Type             string `json:"type"`
}

type Value struct {
	DateTime time.Time `json:"datetime"`
	Open     float64   `json:"open"`
	Close    float64   `json:"low"`
	Low      float64   `json:"high"`
	High     float64   `json:"close"`
	Volume   float64   `json:"volume"`
}

func (v *Value) UnmarshalJSON(b []byte) error {
	type RawValue struct {
		DateTime string `json:"datetime"`
		Open     string `json:"open"`
		Close    string `json:"low"`
		Low      string `json:"high"`
		High     string `json:"close"`
		Volume   string `json:"volume"`
	}

	var rawValue RawValue
	if err := json.Unmarshal(b, &rawValue); err != nil {
		return err
	}

	open, err := strconv.ParseFloat(rawValue.Open, 64)
	if err != nil {
		return errors.Wrap(err, "failed to parse value open into float")
	}
	close, err := strconv.ParseFloat(rawValue.Close, 64)
	if err != nil {
		return errors.Wrap(err, "failed to parse value close into float")
	}

	high, err := strconv.ParseFloat(rawValue.High, 64)
	if err != nil {
		return errors.Wrap(err, "failed to parse value high into float")
	}

	low, err := strconv.ParseFloat(rawValue.Low, 64)
	if err != nil {
		return errors.Wrap(err, "failed to parse value low into float")
	}

	var volume float64
	if rawValue.Volume != "" {
		volume, err = strconv.ParseFloat(rawValue.Volume, 64)
		if err != nil {
			return errors.Wrap(err, "failed to parse value volume into float")
		}

	}

	dateTime, err := time.Parse(model.GetTimeFormatFromString(rawValue.DateTime), rawValue.DateTime)
	if err != nil {
		return errors.Wrap(err, "failed to parse value date time into go time")
	}

	v.DateTime = dateTime
	v.Open = open
	v.Low = low
	v.High = high
	v.Close = close
	v.Volume = volume

	return nil
}

// TimeSeriesOptions - options for calling the twelvedata time series endpoint: https://twelvedata.com/docs#time-series
type TimeSeriesOptions struct {
	Exchange   string
	MICCode    string
	Country    string
	SeriesType string
	Type       string
	OutputSize int
	StartDate  *time.Time
	EndDate    *time.Time
}

func (t TimeSeriesOptions) params(u *url.URL, urlValues url.Values) {
	if t.Exchange != "" {
		urlValues.Add("exchange", t.Exchange)
	}

	if t.MICCode != "" {
		urlValues.Add("mic_code", t.MICCode)
	}

	if t.Country != "" {
		urlValues.Add("country", t.Country)
	}

	if t.SeriesType != "" {
		urlValues.Add("series_type", t.SeriesType)
	}

	if t.Type != "" {
		urlValues.Add("type", t.Type)
	}

	if t.OutputSize == 0 {
		urlValues.Add("outputsize", strconv.Itoa(t.OutputSize))
	}

	if t.StartDate != nil {
		urlValues.Add("start_date", t.StartDate.Format(model.TimeFormatMap[model.OneHour]))
	}

	if t.EndDate != nil {
		urlValues.Add("end_date", t.EndDate.Format(model.TimeFormatMap[model.OneHour]))
	}

	u.RawQuery = urlValues.Encode()
}

func (c *client) TimeSeries(symbol string, interval model.Interval, opts TimeSeriesOptions) (TimeSeriesResponse, error) {
	u, err := url.Parse(fmt.Sprintf("%s/time_series", baseURI))
	if err != nil {
		return TimeSeriesResponse{}, errors.Wrapf(err, "failed to parse base URL '%s'", baseURI)
	}

	opts.params(u, url.Values{
		"symbol":   {symbol},
		"interval": {string(interval)},
		"apikey":   {c.apiKey},
	})

	body, err := c.getFn(u)
	if err != nil {
		return TimeSeriesResponse{}, err
	}

	var response TimeSeriesResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return TimeSeriesResponse{}, err
	}

	return response, nil
}
