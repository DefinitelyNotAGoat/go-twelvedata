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
	stochBody = []byte(`{"meta":{"symbol":"AAPL","interval":"1min","currency":"USD","exchange_timezone":"America/New_York","exchange":"NASDAQ","mic_code":"XNGS","type":"Common Stock","indicator":{"name":"STOCH - Stochastic Oscillator","fast_k_period":14,"slow_k_period":1,"slow_d_period":3,"slow_kma_type":"SMA","slow_dma_type":"SMA"}},"values":[{"datetime":"2023-08-24 10:47:00","slow_k":"37.51669","slow_d":"21.95264"},{"datetime":"2023-08-24 10:46:00","slow_k":"18.75132","slow_d":"14.15896"},{"datetime":"2023-08-24 10:45:00","slow_k":"9.58990","slow_d":"8.56535"},{"datetime":"2023-08-24 10:44:00","slow_k":"14.13565","slow_d":"7.11201"},{"datetime":"2023-08-24 10:43:00","slow_k":"1.97049","slow_d":"4.04048"},{"datetime":"2023-08-24 10:42:00","slow_k":"5.22988","slow_d":"4.94688"},{"datetime":"2023-08-24 10:41:00","slow_k":"4.92107","slow_d":"8.20346"},{"datetime":"2023-08-24 10:40:00","slow_k":"4.68970","slow_d":"14.34097"},{"datetime":"2023-08-24 10:39:00","slow_k":"14.99962","slow_d":"20.99018"},{"datetime":"2023-08-24 10:38:00","slow_k":"23.33359","slow_d":"22.20129"},{"datetime":"2023-08-24 10:37:00","slow_k":"24.63733","slow_d":"15.55974"},{"datetime":"2023-08-24 10:36:00","slow_k":"18.63295","slow_d":"10.66187"},{"datetime":"2023-08-24 10:35:00","slow_k":"3.40893","slow_d":"4.64046"},{"datetime":"2023-08-24 10:34:00","slow_k":"9.94371","slow_d":"5.82963"},{"datetime":"2023-08-24 10:33:00","slow_k":"0.56873","slow_d":"8.95435"},{"datetime":"2023-08-24 10:32:00","slow_k":"6.97646","slow_d":"11.76854"},{"datetime":"2023-08-24 10:31:00","slow_k":"19.31787","slow_d":"9.91617"},{"datetime":"2023-08-24 10:30:00","slow_k":"9.01131","slow_d":"6.35785"},{"datetime":"2023-08-24 10:29:00","slow_k":"1.41934","slow_d":"10.96717"},{"datetime":"2023-08-24 10:28:00","slow_k":"8.64291","slow_d":"13.39966"},{"datetime":"2023-08-24 10:27:00","slow_k":"22.83927","slow_d":"11.56577"},{"datetime":"2023-08-24 10:26:00","slow_k":"8.71681","slow_d":"4.52408"},{"datetime":"2023-08-24 10:25:00","slow_k":"3.14123","slow_d":"3.74605"},{"datetime":"2023-08-24 10:24:00","slow_k":"1.71422","slow_d":"4.40833"},{"datetime":"2023-08-24 10:23:00","slow_k":"6.38270","slow_d":"7.26795"},{"datetime":"2023-08-24 10:22:00","slow_k":"5.12807","slow_d":"9.38429"},{"datetime":"2023-08-24 10:21:00","slow_k":"10.29306","slow_d":"8.06747"},{"datetime":"2023-08-24 10:20:00","slow_k":"12.73174","slow_d":"7.11299"},{"datetime":"2023-08-24 10:19:00","slow_k":"1.17761","slow_d":"5.15494"},{"datetime":"2023-08-24 10:18:00","slow_k":"7.42963","slow_d":"5.79171"}],"status":"ok"}`)
)

func TestIntegrationStochastic(t *testing.T) {
	client := New(
		os.Getenv("TWELVEDATA_API_KEY"),
		http.DefaultClient,
	)

	_, err := client.Stochastic("AAPL", model.OneHour, StochasticOptions{
		IndicatorOptions: IndicatorOptions{
			IncludeOHLC: true,
		},
	})
	if err != nil {
		t.Log("Failed to make Stochastic request: ", err.Error())
		t.Fail()
	}
}

func TestUnitStochastic(t *testing.T) {
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
					return stochBody, nil
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

			_, err := client.Stochastic("AAPL", model.OneDay, StochasticOptions{})
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
