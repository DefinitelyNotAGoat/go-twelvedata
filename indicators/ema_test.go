package indicators

import (
	"errors"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/DefinitelyNotAGoat/twelvedata/model"
	"github.com/stretchr/testify/assert"
)

var (
	emaBody = []byte(`{"meta":{"symbol":"AAPL","interval":"1min","currency":"USD","exchange_timezone":"America/New_York","exchange":"NASDAQ","mic_code":"XNGS","type":"Common Stock","indicator":{"name":"EMA - Exponential Moving Average","series_type":"close","time_period":9}},"values":[{"datetime":"2023-08-24 10:54:00","ema":"178.08526"},{"datetime":"2023-08-24 10:53:00","ema":"178.05465"},{"datetime":"2023-08-24 10:52:00","ema":"177.99832"},{"datetime":"2023-08-24 10:51:00","ema":"177.97289"},{"datetime":"2023-08-24 10:50:00","ema":"177.93672"},{"datetime":"2023-08-24 10:49:00","ema":"177.91090"},{"datetime":"2023-08-24 10:48:00","ema":"177.92362"},{"datetime":"2023-08-24 10:47:00","ema":"177.91828"},{"datetime":"2023-08-24 10:46:00","ema":"177.91782"},{"datetime":"2023-08-24 10:45:00","ema":"177.94353"},{"datetime":"2023-08-24 10:44:00","ema":"177.99191"},{"datetime":"2023-08-24 10:43:00","ema":"178.04266"},{"datetime":"2023-08-24 10:42:00","ema":"178.12582"},{"datetime":"2023-08-24 10:41:00","ema":"178.19478"},{"datetime":"2023-08-24 10:40:00","ema":"178.22597"},{"datetime":"2023-08-24 10:39:00","ema":"178.26496"},{"datetime":"2023-08-24 10:38:00","ema":"178.28870"},{"datetime":"2023-08-24 10:37:00","ema":"178.30588"},{"datetime":"2023-08-24 10:36:00","ema":"178.31984"},{"datetime":"2023-08-24 10:35:00","ema":"178.34231"},{"datetime":"2023-08-24 10:34:00","ema":"178.40038"},{"datetime":"2023-08-24 10:33:00","ema":"178.44561"},{"datetime":"2023-08-24 10:32:00","ema":"178.49826"},{"datetime":"2023-08-24 10:31:00","ema":"178.54532"},{"datetime":"2023-08-24 10:30:00","ema":"178.57665"},{"datetime":"2023-08-24 10:29:00","ema":"178.63849"},{"datetime":"2023-08-24 10:28:00","ema":"178.69761"},{"datetime":"2023-08-24 10:27:00","ema":"178.74202"},{"datetime":"2023-08-24 10:26:00","ema":"178.76877"},{"datetime":"2023-08-24 10:25:00","ema":"178.82721"}],"status":"ok"}`)
)

func TestIntegrationEMA(t *testing.T) {
	client := New(
		os.Getenv("TWELVEDATA_API_KEY"),
		http.DefaultClient,
	)

	_, err := client.EMA("AAPL", model.OneHour, EMAOptions{
		IndicatorOptions: IndicatorOptions{
			IncludeOHLC: true,
		},
	})
	if err != nil {
		t.Log("Failed to make EMA request: ", err.Error())
		t.Fail()
	}
}

func TestUnitEMA(t *testing.T) {
	type input struct {
		getFn getFn
	}

	type want struct {
		err      bool
		contains string
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"handles failure to get",
			input{
				getFn: func(u *url.URL) ([]byte, error) {
					return nil, errors.New("failed to get")
				},
			},
			want{
				err:      true,
				contains: "failed to get",
			},
		},
		{
			"is successful",
			input{
				getFn: func(u *url.URL) ([]byte, error) {
					return emaBody, nil
				},
			},
			want{
				err: false,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			client := client{
				c:     http.DefaultClient,
				getFn: tt.input.getFn,
			}

			_, err := client.EMA("AAPL", model.OneDay, EMAOptions{})
			if tt.want.err {
				if assert.NotNil(t, err) {
					assert.Contains(t, err.Error(), tt.want.contains)
				}
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
