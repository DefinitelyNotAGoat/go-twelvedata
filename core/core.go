package core

import (
	"net/http"
	"net/url"

	httpt "github.com/DefinitelyNotAGoat/twelvedata/http"
	"github.com/DefinitelyNotAGoat/twelvedata/model"
)

const (
	baseURI = "https://api.twelvedata.com"
)

type getFn func(u *url.URL) ([]byte, error)

// Client - Exposes an interface to interact with Twelvedata's core API: https://twelvedata.com/docs#core-data
type Client interface {
	TimeSeries(symbol string, interval model.Interval, opts TimeSeriesOptions) (TimeSeriesResponse, error)
	MarketMovers(opts MarketMoversOptions) (MarketMoversResponse, error)
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
