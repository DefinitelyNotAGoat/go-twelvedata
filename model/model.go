package model

import "strings"

type Interval string

const (
	//OneMin        Interval = "1min"
	//FiveMin       Interval = "5min"
	//FifteenMin    Interval = "15min"
	//ThirtyMin     Interval = "30min"
	//FourtyFiveMin Interval = "45min"
	OneHour Interval = "1h"
	//TwoHour       Interval = "2h"
	FourHour Interval = "4h"
	OneDay   Interval = "1day"
	OneWeek  Interval = "1week"
	OneMonth Interval = "1month"
)

var (
	TimeFormatMap map[Interval]string = map[Interval]string{
		OneHour:  "2006-01-02 15:04:05",
		FourHour: "2006-01-02 15:04:05",
		OneDay:   "2006-01-02",
		OneWeek:  "2006-01-02",
		OneMonth: "2006-01-02",
	}
)

func GetTimeFormatFromString(str string) string {
	strs := strings.Split(str, ":")
	if len(strs) == 1 {
		return TimeFormatMap[OneDay]
	}

	return TimeFormatMap[OneHour]
}

type Meta struct {
	Symbol           string `json:"symbol"`
	Interval         string `json:"interval"`
	Currency         string `json:"currency"`
	ExchangeTimezone string `json:"exchange_timezone"`
	Exchange         string `json:"exchange"`
	MicCode          string `json:"mic_code"`
	Type             string `json:"type"`
}
