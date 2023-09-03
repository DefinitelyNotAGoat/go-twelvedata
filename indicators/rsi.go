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

// RSIIndicator - the Indiciator value for IndicatorMeta specific for RSI
type RSIIndicator struct {
	Name       string `json:"name"`
	SeriesType string `json:"series_type"`
	TimePeriod int    `json:"time_period"`
}

// RSIValue - the Indicator value for IndicatorResponse specific for RSI
type RSIValue struct {
	Datetime time.Time `json:"datetime"`
	Rsi      float64   `json:"rsi"`
}

// UnmarshalJSON - unmarshal's RSIValue to a more consumable type
func (r *RSIValue) UnmarshalJSON(v []byte) error {
	type rsiValue struct {
		Datetime string `json:"datetime"`
		Rsi      string `json:"rsi"`
	}

	var value rsiValue
	if err := json.Unmarshal(v, &value); err != nil {
		return err
	}

	rsi, err := strconv.ParseFloat(value.Rsi, 64)
	if err != nil {
		return errors.Wrap(err, "failed to parse value ema into float")
	}

	dateTime, err := time.Parse(model.GetTimeFormatFromString(value.Datetime), value.Datetime)
	if err != nil {
		return errors.Wrap(err, "failed to parse value date time into go time")
	}

	r.Rsi = rsi
	r.Datetime = dateTime

	return nil
}

// RSIOptions - options for calling the twelvedata rsi endpoint: https://twelvedata.com/docs#rsi
type RSIOptions struct {
	IndicatorOptions
	TimePeriod int
}

func (r RSIOptions) params(u *url.URL, urlValues url.Values) {
	urlValues = r.IndicatorOptions.params(u, urlValues)

	if r.TimePeriod > 0 {
		urlValues.Add("time_period", strconv.Itoa(r.TimePeriod))
	}

	u.RawQuery = urlValues.Encode()
}

func rsi[Response IndicatorResponse[RSIValue, RSIIndicator]](symbol string, interval model.Interval, apiKey string, getFn getFn, opts RSIOptions) (Response, error) {
	u, err := url.Parse(fmt.Sprintf("%s/rsi", baseURI))
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
