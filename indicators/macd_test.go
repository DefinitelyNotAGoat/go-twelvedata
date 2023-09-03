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
	macdBody = []byte(`{"meta":{"symbol":"AAPL","interval":"1min","currency":"USD","exchange_timezone":"America/New_York","exchange":"NASDAQ","mic_code":"XNGS","type":"Common Stock","indicator":{"name":"MACD - Moving Average Convergence Divergence","series_type":"close","fast_period":12,"slow_period":26,"signal_period":9}},"values":[{"datetime":"2023-08-23 15:59:00","macd":"-0.08199","macd_signal":"-0.05782","macd_hist":"-0.02417"},{"datetime":"2023-08-23 15:58:00","macd":"-0.08514","macd_signal":"-0.05178","macd_hist":"-0.03336"},{"datetime":"2023-08-23 15:57:00","macd":"-0.08028","macd_signal":"-0.04344","macd_hist":"-0.03685"},{"datetime":"2023-08-23 15:56:00","macd":"-0.07649","macd_signal":"-0.03422","macd_hist":"-0.04226"},{"datetime":"2023-08-23 15:55:00","macd":"-0.06900","macd_signal":"-0.02366","macd_hist":"-0.04534"},{"datetime":"2023-08-23 15:54:00","macd":"-0.05538","macd_signal":"-0.01232","macd_hist":"-0.04306"},{"datetime":"2023-08-23 15:53:00","macd":"-0.04229","macd_signal":"-0.00156","macd_hist":"-0.04073"},{"datetime":"2023-08-23 15:52:00","macd":"-0.03005","macd_signal":"0.00863","macd_hist":"-0.03868"},{"datetime":"2023-08-23 15:51:00","macd":"-0.01196","macd_signal":"0.01830","macd_hist":"-0.03026"},{"datetime":"2023-08-23 15:50:00","macd":"-0.00041","macd_signal":"0.02586","macd_hist":"-0.02627"},{"datetime":"2023-08-23 15:49:00","macd":"0.01115","macd_signal":"0.03243","macd_hist":"-0.02128"},{"datetime":"2023-08-23 15:48:00","macd":"0.01875","macd_signal":"0.03775","macd_hist":"-0.01900"},{"datetime":"2023-08-23 15:47:00","macd":"0.01851","macd_signal":"0.04250","macd_hist":"-0.02399"},{"datetime":"2023-08-23 15:46:00","macd":"0.02405","macd_signal":"0.04850","macd_hist":"-0.02446"},{"datetime":"2023-08-23 15:45:00","macd":"0.03534","macd_signal":"0.05462","macd_hist":"-0.01927"},{"datetime":"2023-08-23 15:44:00","macd":"0.04874","macd_signal":"0.05943","macd_hist":"-0.01070"},{"datetime":"2023-08-23 15:43:00","macd":"0.05443","macd_signal":"0.06211","macd_hist":"-0.00768"},{"datetime":"2023-08-23 15:42:00","macd":"0.05327","macd_signal":"0.06403","macd_hist":"-0.01076"},{"datetime":"2023-08-23 15:41:00","macd":"0.05509","macd_signal":"0.06672","macd_hist":"-0.01162"},{"datetime":"2023-08-23 15:40:00","macd":"0.05763","macd_signal":"0.06962","macd_hist":"-0.01199"},{"datetime":"2023-08-23 15:39:00","macd":"0.05800","macd_signal":"0.07262","macd_hist":"-0.01462"},{"datetime":"2023-08-23 15:38:00","macd":"0.06424","macd_signal":"0.07628","macd_hist":"-0.01204"},{"datetime":"2023-08-23 15:37:00","macd":"0.06618","macd_signal":"0.07929","macd_hist":"-0.01310"},{"datetime":"2023-08-23 15:36:00","macd":"0.07165","macd_signal":"0.08256","macd_hist":"-0.01091"},{"datetime":"2023-08-23 15:35:00","macd":"0.07588","macd_signal":"0.08529","macd_hist":"-0.00942"},{"datetime":"2023-08-23 15:34:00","macd":"0.07996","macd_signal":"0.08765","macd_hist":"-0.00768"},{"datetime":"2023-08-23 15:33:00","macd":"0.08401","macd_signal":"0.08957","macd_hist":"-0.00556"},{"datetime":"2023-08-23 15:32:00","macd":"0.09011","macd_signal":"0.09096","macd_hist":"-0.00085"},{"datetime":"2023-08-23 15:31:00","macd":"0.09799","macd_signal":"0.09117","macd_hist":"0.00683"},{"datetime":"2023-08-23 15:30:00","macd":"0.09594","macd_signal":"0.08946","macd_hist":"0.00648"}],"status":"ok"}`)
)

func TestIntegrationMACD(t *testing.T) {
	client := New(
		os.Getenv("TWELVEDATA_API_KEY"),
		http.DefaultClient,
	)

	_, err := client.MACD("AAPL", model.OneHour, MACDOptions{
		IndicatorOptions: IndicatorOptions{
			IncludeOHLC: true,
		},
	})
	if err != nil {
		t.Log("Failed to make MACD request: ", err.Error())
		t.Fail()
	}
}

func TestUnitMACD(t *testing.T) {
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
					return macdBody, nil
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

			_, err := client.MACD("AAPL", model.OneDay, MACDOptions{})
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
