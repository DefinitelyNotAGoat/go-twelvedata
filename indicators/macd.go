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

// MACDIndicator - the Indiciator value for IndicatorMeta specific for MACD
type MACDIndicator struct {
	Name         string `json:"name"`
	SeriesType   string `json:"series_type"`
	FastPeriod   int    `json:"fast_period"`
	SlowPeriod   int    `json:"slow_period"`
	SignalPeriod int    `json:"signal_period"`
}

// MACDValue - the Indicator value for IndicatorResponse specific for MACD
type MACDValue struct {
	Datetime   time.Time `json:"datetime"`
	Macd       float64   `json:"macd"`
	MacdSignal float64   `json:"macd_signal"`
	MacdHist   float64   `json:"macd_hist"`
}

// UnmarshalJSON - unmarshal's MACDValue to a more consumable type
func (m *MACDValue) UnmarshalJSON(v []byte) error {
	type Value struct {
		Datetime   string `json:"datetime"`
		Macd       string `json:"macd"`
		MacdSignal string `json:"macd_signal"`
		MacdHist   string `json:"macd_hist"`
	}

	var value Value
	if err := json.Unmarshal(v, &value); err != nil {
		return err
	}

	macd, err := strconv.ParseFloat(value.Macd, 64)
	if err != nil {
		return errors.Wrap(err, "failed to parse value macd into float")
	}
	signal, err := strconv.ParseFloat(value.MacdSignal, 64)
	if err != nil {
		return errors.Wrap(err, "failed to parse value signal into float")
	}

	historgam, err := strconv.ParseFloat(value.MacdHist, 64)
	if err != nil {
		return errors.Wrap(err, "failed to parse value histogram into float")
	}

	dateTime, err := time.Parse(model.GetTimeFormatFromString(value.Datetime), value.Datetime)
	if err != nil {
		return errors.Wrap(err, "failed to parse value date time into go time")
	}

	m.Datetime = dateTime
	m.Macd = macd
	m.MacdSignal = signal
	m.MacdHist = historgam

	return nil
}

// MACDOptions - options for calling the twelvedata macd endpoint: https://twelvedata.com/docs#macd
type MACDOptions struct {
	IndicatorOptions
	FastPeriod   int
	SignalPeriod int
	SlowPeriod   int
}

func (m MACDOptions) params(u *url.URL, urlValues url.Values) {
	urlValues = m.IndicatorOptions.params(u, urlValues)

	if m.FastPeriod != 0 {
		urlValues.Add("fast_period", strconv.Itoa(m.FastPeriod))
	}

	if m.SignalPeriod != 0 {
		urlValues.Add("signal_period", strconv.Itoa(m.SignalPeriod))
	}

	if m.SlowPeriod != 0 {
		urlValues.Add("slow_period", strconv.Itoa(m.SlowPeriod))
	}

	u.RawQuery = urlValues.Encode()
}

func macd[Response IndicatorResponse[MACDValue, MACDIndicator]](symbol string, interval model.Interval, apiKey string, getFn getFn, opts MACDOptions) (Response, error) {
	u, err := url.Parse(fmt.Sprintf("%s/macd", baseURI))
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
