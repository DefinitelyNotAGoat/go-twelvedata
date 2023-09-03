package core

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
	timeSeriesBody = []byte(`{"meta":{"symbol":"AAPL","interval":"1min","currency":"USD","exchange_timezone":"America/New_York","exchange":"NASDAQ","mic_code":"XNGS","type":"Common Stock"},"values":[{"datetime":"2023-08-24 11:09:00","open":"177.91010","high":"178.02000","low":"177.91000","close":"178.00121","volume":"293189"},{"datetime":"2023-08-24 11:08:00","open":"177.75500","high":"177.88499","low":"177.73000","close":"177.88000","volume":"143249"},{"datetime":"2023-08-24 11:07:00","open":"177.73840","high":"177.81000","low":"177.69501","close":"177.74500","volume":"57897"},{"datetime":"2023-08-24 11:06:00","open":"177.87000","high":"177.89090","low":"177.66499","close":"177.74001","volume":"119186"},{"datetime":"2023-08-24 11:05:00","open":"177.93010","high":"178.05960","low":"177.86031","close":"177.88000","volume":"143464"},{"datetime":"2023-08-24 11:04:00","open":"177.78000","high":"177.99001","low":"177.75000","close":"177.94000","volume":"95722"},{"datetime":"2023-08-24 11:03:00","open":"177.85809","high":"177.87660","low":"177.78000","close":"177.78000","volume":"65281"},{"datetime":"2023-08-24 11:02:00","open":"177.80499","high":"177.86501","low":"177.74001","close":"177.85899","volume":"62373"},{"datetime":"2023-08-24 11:01:00","open":"177.75000","high":"177.89999","low":"177.69400","close":"177.80499","volume":"127906"},{"datetime":"2023-08-24 11:00:00","open":"177.91000","high":"177.91000","low":"177.64999","close":"177.75999","volume":"170797"},{"datetime":"2023-08-24 10:59:00","open":"178.11000","high":"178.11000","low":"177.88139","close":"177.89000","volume":"104104"},{"datetime":"2023-08-24 10:58:00","open":"178.14169","high":"178.27000","low":"178.10060","close":"178.11501","volume":"97989"},{"datetime":"2023-08-24 10:57:00","open":"178.26801","high":"178.28500","low":"178.18500","close":"178.25900","volume":"130845"},{"datetime":"2023-08-24 10:56:00","open":"178.27000","high":"178.31010","low":"178.25980","close":"178.27499","volume":"75952"},{"datetime":"2023-08-24 10:55:00","open":"178.20500","high":"178.28900","low":"178.19000","close":"178.26500","volume":"60764"},{"datetime":"2023-08-24 10:54:00","open":"178.27901","high":"178.31990","low":"178.19000","close":"178.20770","volume":"105089"},{"datetime":"2023-08-24 10:53:00","open":"178.09000","high":"178.28000","low":"178.07001","close":"178.28000","volume":"127884"},{"datetime":"2023-08-24 10:52:00","open":"178.10001","high":"178.11990","low":"177.99989","close":"178.10001","volume":"105334"},{"datetime":"2023-08-24 10:51:00","open":"178.03999","high":"178.15500","low":"178.03000","close":"178.11760","volume":"143380"},{"datetime":"2023-08-24 10:50:00","open":"177.86000","high":"178.03999","low":"177.84000","close":"178.03999","volume":"103858"},{"datetime":"2023-08-24 10:49:00","open":"177.94000","high":"177.98351","low":"177.84000","close":"177.86000","volume":"141962"},{"datetime":"2023-08-24 10:48:00","open":"177.92310","high":"177.94501","low":"177.88080","close":"177.94501","volume":"97955"},{"datetime":"2023-08-24 10:47:00","open":"177.82500","high":"177.94501","low":"177.81500","close":"177.92010","volume":"120554"},{"datetime":"2023-08-24 10:46:00","open":"177.73790","high":"177.82001","low":"177.71500","close":"177.81500","volume":"101057"},{"datetime":"2023-08-24 10:45:00","open":"177.79111","high":"177.83000","low":"177.67999","close":"177.75000","volume":"130006"},{"datetime":"2023-08-24 10:44:00","open":"177.71500","high":"177.82001","low":"177.68500","close":"177.78889","volume":"208785"},{"datetime":"2023-08-24 10:43:00","open":"177.85500","high":"177.89999","low":"177.69501","close":"177.71001","volume":"184111"},{"datetime":"2023-08-24 10:42:00","open":"178.06000","high":"178.07001","low":"177.81000","close":"177.85001","volume":"293147"},{"datetime":"2023-08-24 10:41:00","open":"178.09000","high":"178.12000","low":"178.05000","close":"178.07001","volume":"104754"},{"datetime":"2023-08-24 10:40:00","open":"178.16000","high":"178.22060","low":"178.03999","close":"178.07001","volume":"175143"}],"status":"ok"}`)
)

func TestIntegrationTimeSeries(t *testing.T) {
	client := New(
		os.Getenv("TWELVEDATA_API_KEY"),
		http.DefaultClient,
	)

	_, err := client.TimeSeries("AAPL", model.OneHour, TimeSeriesOptions{})
	if err != nil {
		t.Log("Failed to make TimeSeries request: ", err.Error())
		t.Fail()
	}
}

func TestUnitTimeSeries(t *testing.T) {
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
					return timeSeriesBody, nil
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

			_, err := client.TimeSeries("AAPL", model.OneDay, TimeSeriesOptions{})
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
