package twelvedata

import (
	"net/http"

	"github.com/DefinitelyNotAGoat/twelvedata/core"
	"github.com/DefinitelyNotAGoat/twelvedata/indicators"
)

// Client - a general wrapper that encomposes all TwelveData API groups
type Client struct {
	CoreData            core.Client
	TechnicalIndicators indicators.Client
}

// New - returns a new TwelveData Client
func New(apiKey string, client *http.Client) Client {
	return Client{
		CoreData:            core.New(apiKey, client),
		TechnicalIndicators: indicators.New(apiKey, client),
	}
}
