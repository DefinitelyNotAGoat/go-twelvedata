package core

import (
	"errors"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	marketMoversBody = []byte(`{"values":[{"symbol":"SQL","name":"SeqLL Inc.","exchange":"NASDAQ","mic_code":"XNCM","datetime":"2023-08-31 10:11:42","last":17.3261,"high":21.9799,"low":15,"volume":60185,"change":17.2002,"percent_change":3823.96621},{"symbol":"EXPR","name":"Express Inc","exchange":"NYSE","mic_code":"XNYS","datetime":"2023-08-31 10:13:05","last":9.71,"high":9.7,"low":9.02,"volume":56634,"change":8.612,"percent_change":1695.27559},{"symbol":"SPIR","name":"Spire Corp.","exchange":"NYSE","mic_code":"XNYS","datetime":"2023-08-31 10:09:18","last":4.48,"high":4.7,"low":4.35,"volume":99688,"change":3.9556,"percent_change":713.49206},{"symbol":"ICCT","name":"iCoreConnect, Inc.","exchange":"NASDAQ","mic_code":"XNCM","datetime":"2023-08-31 10:12:54","last":14.81,"high":20.7,"low":13.63,"volume":8628114,"change":1.0785,"percent_change":563.18537},{"symbol":"LFLYW","name":"Leafly Holdings, Inc.","exchange":"NASDAQ","mic_code":"XNMS","datetime":"2023-08-31 10:04:00","last":0.04,"high":0.04,"low":0.035,"volume":0,"change":0.0195,"percent_change":95.12195},{"symbol":"AGIL","name":"AgileThought, Inc.","exchange":"NASDAQ","mic_code":"XNCM","datetime":"2023-08-31 10:12:26","last":0.2026,"high":0.2292,"low":0.2,"volume":3299943,"change":0.119,"percent_change":70.41419},{"symbol":"NCNC","name":"Noco-noco Inc","exchange":"NASDAQ","mic_code":"XNCM","datetime":"2023-08-31 10:07:00","last":2.6199,"high":2.75,"low":2.1,"volume":26768539,"change":1.0209,"percent_change":63.84616},{"symbol":"ACER","name":"Acer Therapeutics Inc","exchange":"NASDAQ","mic_code":"XNCM","datetime":"2023-08-31 10:13:28","last":1.57,"high":1.63,"low":0.9131,"volume":55615519,"change":0.3795,"percent_change":59.29688},{"symbol":"SATX","name":"SatixFy Communications Ltd.","exchange":"NYSE","mic_code":"XASE","datetime":"2023-08-31 10:12:55","last":0.6,"high":0.7188,"low":0.5855,"volume":13682016,"change":0.2345,"percent_change":51.25683},{"symbol":"KACLR","name":"Kairous Acquisition Corp. Limited","exchange":"NASDAQ","mic_code":"XNMS","datetime":"2023-08-31 09:48:00","last":0.119,"high":0.119,"low":0.119,"volume":0,"change":0.035,"percent_change":41.66667},{"symbol":"ZVZZT","name":"","exchange":"NASDAQ","mic_code":"XNGS","datetime":"2023-08-31 09:34:00","last":25.27,"high":25.27,"low":25.08,"volume":0,"change":6.68,"percent_change":35.9333},{"symbol":"BBIG","name":"Vinco Ventures, Inc.","exchange":"NASDAQ","mic_code":"XNCM","datetime":"2023-08-31 09:51:00","last":0.08,"high":0.18,"low":0.08,"volume":724,"change":0.02,"percent_change":33.33333},{"symbol":"ELMS","name":"Electric Last Mile Solutions, Inc.","exchange":"NASDAQ","mic_code":"XNGS","datetime":"2023-08-31 09:49:00","last":0.09,"high":0.09,"low":0.07,"volume":0,"change":0.02,"percent_change":28.57143},{"symbol":"EDNC","name":"Endurance Acquisition Corp.","exchange":"NASDAQ","mic_code":"XNCM","datetime":"2023-08-31 10:12:00","last":0.6,"high":0.719,"low":0.59,"volume":0,"change":0.131,"percent_change":27.93177},{"symbol":"TIO","name":"Tingo Group, Inc. - Common Stock","exchange":"NASDAQ","mic_code":"XNCM","datetime":"2023-08-31 10:13:01","last":1.47,"high":1.835,"low":1.43,"volume":40861858,"change":0.2253,"percent_change":26.51525},{"symbol":"ALLR","name":"Allarity Therapeutics, Inc.","exchange":"NASDAQ","mic_code":"XNCM","datetime":"2023-08-31 10:12:45","last":1.8013,"high":1.91,"low":1.69,"volume":3886294,"change":0.37,"percent_change":22.28916},{"symbol":"FFIEW","name":"Faraday Future Intelligent Electric Inc.","exchange":"NASDAQ","mic_code":"XNMS","datetime":"2023-08-31 09:37:18","last":0.035,"high":0.035,"low":0.035,"volume":0,"change":0.0085,"percent_change":20.73171},{"symbol":"BRSHW","name":"Bruush Oral Care Inc.","exchange":"NASDAQ","mic_code":"XNCM","datetime":"2023-08-31 09:30:00","last":0.035,"high":0.035,"low":0.035,"volume":0,"change":0.006,"percent_change":20.68966},{"symbol":"XPER","name":"Xperi Corp","exchange":"NYSE","mic_code":"XNYS","datetime":"2023-08-31 10:11:55","last":10.17,"high":12.16,"low":12,"volume":3642,"change":2.03,"percent_change":20.11893},{"symbol":"ATNFW","name":"CannBioRx Life Sciences Corp.","exchange":"NASDAQ","mic_code":"XNCM","datetime":"2023-08-31 10:12:00","last":0.012,"high":0.01,"low":0.01,"volume":0,"change":0.002,"percent_change":20},{"symbol":"VERY","name":"Vericity Inc","exchange":"NASDAQ","mic_code":"XNCM","datetime":"2023-08-31 09:53:13","last":5.2,"high":5.2,"low":4.79,"volume":0,"change":0.8548,"percent_change":17.096},{"symbol":"NLSP","name":"NLS Pharmaceutics AG","exchange":"NASDAQ","mic_code":"XNCM","datetime":"2023-08-31 10:12:30","last":0.9773,"high":1.1462,"low":0.93,"volume":383222,"change":0.16,"percent_change":16},{"symbol":"PROCW","name":"Procaps Group, S.A.","exchange":"NASDAQ","mic_code":"XNMS","datetime":"2023-08-31 09:41:00","last":0.1538,"high":0.1538,"low":0.1538,"volume":0,"change":0.0212,"percent_change":15.98793},{"symbol":"HLTH","name":"Nobilis Health Corp","exchange":"NASDAQ","mic_code":"XNGS","datetime":"2023-08-31 10:10:00","last":0.52,"high":0.52,"low":0.4418,"volume":109663,"change":0.0697,"percent_change":15.83012},{"symbol":"SFIX","name":"Stitch Fix Inc","exchange":"NASDAQ","mic_code":"XNGS","datetime":"2023-08-31 10:13:03","last":4.41,"high":4.43,"low":4.33,"volume":144350,"change":0.59,"percent_change":15.4047}],"status":"ok"}`)
)

func TestIntegrationMarketMovers(t *testing.T) {
	client := New(
		os.Getenv("TWELVEDATA_API_KEY"),
		http.DefaultClient,
	)

	_, err := client.MarketMovers(MarketMoversOptions{})
	if err != nil {
		t.Log("Failed to make MarketMovers request: ", err.Error())
		t.Fail()
	}
}

func TestUnitMarketMovers(t *testing.T) {
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
					return marketMoversBody, nil
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

			_, err := client.MarketMovers(MarketMoversOptions{})
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
