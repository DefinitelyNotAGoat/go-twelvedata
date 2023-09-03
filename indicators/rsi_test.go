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
	rsiBody = []byte(`{"meta":{"symbol":"AAPL","interval":"1min","currency":"USD","exchange_timezone":"America/New_York","exchange":"NASDAQ","mic_code":"XNGS","type":"Common Stock","indicator":{"name":"RSI - Relative Strength Index","series_type":"close","time_period":14}},"values":[{"datetime":"2023-08-24 10:33:00","rsi":"22.89856"},{"datetime":"2023-08-24 10:32:00","rsi":"24.37532"},{"datetime":"2023-08-24 10:31:00","rsi":"26.72250"},{"datetime":"2023-08-24 10:30:00","rsi":"20.89012"},{"datetime":"2023-08-24 10:29:00","rsi":"22.20548"},{"datetime":"2023-08-24 10:28:00","rsi":"24.53422"},{"datetime":"2023-08-24 10:27:00","rsi":"27.10630"},{"datetime":"2023-08-24 10:26:00","rsi":"20.36523"},{"datetime":"2023-08-24 10:25:00","rsi":"18.61729"},{"datetime":"2023-08-24 10:24:00","rsi":"19.65871"},{"datetime":"2023-08-24 10:23:00","rsi":"23.39805"},{"datetime":"2023-08-24 10:22:00","rsi":"25.78980"},{"datetime":"2023-08-24 10:21:00","rsi":"26.75706"},{"datetime":"2023-08-24 10:20:00","rsi":"29.03231"},{"datetime":"2023-08-24 10:19:00","rsi":"25.77707"},{"datetime":"2023-08-24 10:18:00","rsi":"28.17778"},{"datetime":"2023-08-24 10:17:00","rsi":"27.86554"},{"datetime":"2023-08-24 10:16:00","rsi":"28.55688"},{"datetime":"2023-08-24 10:15:00","rsi":"31.58761"},{"datetime":"2023-08-24 10:14:00","rsi":"31.58761"},{"datetime":"2023-08-24 10:13:00","rsi":"31.58761"},{"datetime":"2023-08-24 10:12:00","rsi":"37.12745"},{"datetime":"2023-08-24 10:11:00","rsi":"28.68855"},{"datetime":"2023-08-24 10:10:00","rsi":"23.93329"},{"datetime":"2023-08-24 10:09:00","rsi":"25.51306"},{"datetime":"2023-08-24 10:08:00","rsi":"26.31949"},{"datetime":"2023-08-24 10:07:00","rsi":"29.57550"},{"datetime":"2023-08-24 10:06:00","rsi":"31.76464"},{"datetime":"2023-08-24 10:05:00","rsi":"37.32269"},{"datetime":"2023-08-24 10:04:00","rsi":"39.11557"}],"status":"ok"}`)
)

func TestIntegrationRSI(t *testing.T) {
	client := New(
		os.Getenv("TWELVEDATA_API_KEY"),
		http.DefaultClient,
	)

	_, err := client.RSI("AAPL", model.OneHour, RSIOptions{
		IndicatorOptions: IndicatorOptions{
			IncludeOHLC: true,
		},
	})
	if err != nil {
		t.Log("Failed to make RSI request: ", err.Error())
		t.Fail()
	}
}

func TestUnitRSI(t *testing.T) {
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
					return rsiBody, nil
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

			_, err := client.RSI("AAPL", model.OneDay, RSIOptions{})
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
