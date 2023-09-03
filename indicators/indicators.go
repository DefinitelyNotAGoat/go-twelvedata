package indicators

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	httpt "github.com/DefinitelyNotAGoat/twelvedata/http"
	"github.com/DefinitelyNotAGoat/twelvedata/model"
)

const (
	baseURI = "https://api.twelvedata.com"
)

type getFn func(u *url.URL) ([]byte, error)

// Client - Exposes an interface to interact with Twelvedata's technical indicators API: https://twelvedata.com/docs#technical-indicators
type Client interface {
	EMA(symbol string, interval model.Interval, opts EMAOptions) (IndicatorResponse[EMAValue, EMAIndicator], error)
	MACD(symbol string, interval model.Interval, opts MACDOptions) (IndicatorResponse[MACDValue, MACDIndicator], error)
	RSI(symbol string, interval model.Interval, opts RSIOptions) (IndicatorResponse[RSIValue, RSIIndicator], error)
	Stochastic(symbol string, interval model.Interval, opts StochasticOptions) (IndicatorResponse[StochasticValue, StochasticIndicator], error)
}

type client struct {
	apiKey string
	c      *http.Client
	getFn  getFn
}

// New - returns a new Twelvedata's technical indicators Client
func New(apiKey string, c *http.Client) Client {
	return &client{
		apiKey: apiKey,
		c:      c,
		getFn: func(u *url.URL) ([]byte, error) {
			return httpt.Get(u, c)
		},
	}
}

// EMA - is a generic indicator function for getting the EMA: https://twelvedata.com/docs#ema
func (c *client) EMA(symbol string, interval model.Interval, opts EMAOptions) (IndicatorResponse[EMAValue, EMAIndicator], error) {
	return ema(symbol, interval, c.apiKey, c.getFn, opts)
}

// MACD - is a generic indicator function for getting the MACD: https://twelvedata.com/docs#macd
func (c *client) MACD(symbol string, interval model.Interval, opts MACDOptions) (IndicatorResponse[MACDValue, MACDIndicator], error) {
	return macd(symbol, interval, c.apiKey, c.getFn, opts)
}

// RSI - is a generic indicator function for getting the RSI: https://twelvedata.com/docs#rsi
func (c *client) RSI(symbol string, interval model.Interval, opts RSIOptions) (IndicatorResponse[RSIValue, RSIIndicator], error) {
	return rsi(symbol, interval, c.apiKey, c.getFn, opts)
}

// Stochastic - is a generic indicator function for getting the Stochastic: https://twelvedata.com/docs#stoch
func (c *client) Stochastic(symbol string, interval model.Interval, opts StochasticOptions) (IndicatorResponse[StochasticValue, StochasticIndicator], error) {
	return stochastic(symbol, interval, c.apiKey, c.getFn, opts)
}

// IndicatorResponse - the shared response received from hitting twelvedata's indicator endpoints
type IndicatorResponse[V IndicatorValue, I Indicator] struct {
	Meta   IndicatorMeta[I] `json:"meta"`
	Values []V              `json:"values"`
	Status string           `json:"status"`
}

// IndicatorMeta - a shared substructure of IndicatorResponse
type IndicatorMeta[I Indicator] struct {
	model.Meta
	Indicator I `json:"indicator"`
}

// IndicatorValue - A generic type representing the Values field on the shared IndicatorResponse values
type IndicatorValue interface {
	EMAValue | MACDValue | RSIValue | StochasticValue
}

// Indicator - A generic type respresenting the Indicator field on the shared IndicatorMeta values
type Indicator interface {
	EMAIndicator | MACDIndicator | RSIIndicator | StochasticIndicator
}

// IndicatorOptions - common url query options for all indicator based requests
type IndicatorOptions struct {
	Exchange    string
	MICCode     string
	Country     string
	SeriesType  string
	Type        string
	OutputSize  int
	IncludeOHLC bool
	StartDate   *time.Time
	EndDate     *time.Time
}

func (i IndicatorOptions) params(u *url.URL, urlValues url.Values) url.Values {
	if i.Exchange != "" {
		urlValues.Add("exchange", i.Exchange)
	}

	if i.MICCode != "" {
		urlValues.Add("mic_code", i.MICCode)
	}

	if i.Country != "" {
		urlValues.Add("country", i.Country)
	}

	if i.SeriesType != "" {
		urlValues.Add("series_type", i.SeriesType)
	}

	if i.Type != "" {
		urlValues.Add("type", i.Type)
	}

	if i.OutputSize == 0 {
		urlValues.Add("outputsize", strconv.Itoa(i.OutputSize))
	}

	if i.IncludeOHLC {
		urlValues.Add("include_ohlc", "true")
	}

	if i.StartDate != nil {
		urlValues.Add("start_date", i.StartDate.Format(model.TimeFormatMap[model.OneHour]))
	}

	if i.EndDate != nil {
		urlValues.Add("end_date", i.EndDate.Format(model.TimeFormatMap[model.OneHour]))
	}

	return urlValues
}
