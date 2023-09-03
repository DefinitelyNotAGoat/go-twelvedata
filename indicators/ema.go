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

// EMAIndicator - the Indiciator value for IndicatorMeta specific for EMA
type EMAIndicator struct {
	Name       string `json:"name"`
	SeriesType string `json:"series_type"`
	TimePeriod int    `json:"time_period"`
}

// EMAValue - the Indicator value for IndicatorResponse specific for EMA
type EMAValue struct {
	Datetime time.Time `json:"datetime"`
	Ema      float64   `json:"ema"`
}

// UnmarshalJSON - unmarshal's EMAValue to a more consumable type
func (e *EMAValue) UnmarshalJSON(v []byte) error {
	type emaValue struct {
		Datetime string `json:"datetime"`
		Ema      string `json:"ema"`
	}

	var value emaValue
	if err := json.Unmarshal(v, &value); err != nil {
		return err
	}

	ema, err := strconv.ParseFloat(value.Ema, 64)
	if err != nil {
		return errors.Wrap(err, "failed to parse value ema into float")
	}

	dateTime, err := time.Parse(model.GetTimeFormatFromString(value.Datetime), value.Datetime)
	if err != nil {
		return errors.Wrap(err, "failed to parse value date time into go time")
	}

	e.Ema = ema
	e.Datetime = dateTime

	return nil
}

// EMAOptions - options for calling the twelvedata ema endpoint: https://twelvedata.com/docs#ema
type EMAOptions struct {
	IndicatorOptions
	TimePeriod int
}

func (e EMAOptions) params(u *url.URL, urlValues url.Values) {
	urlValues = e.IndicatorOptions.params(u, urlValues)

	if e.TimePeriod > 0 {
		urlValues.Add("time_period", strconv.Itoa(e.TimePeriod))
	}

	u.RawQuery = urlValues.Encode()
}

func ema[Response IndicatorResponse[EMAValue, EMAIndicator]](symbol string, interval model.Interval, apiKey string, getFn getFn, opts EMAOptions) (Response, error) {
	u, err := url.Parse(fmt.Sprintf("%s/ema", baseURI))
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
