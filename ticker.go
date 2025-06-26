package main

import "encoding/json"

type Ticker map[string]TickerValue

func UnmarshalTicker(data []byte) (Ticker, error) {
	var r Ticker
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Ticker) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type TickerValue struct {
	BuyPrice  string `json:"buy_price"`
	SellPrice string `json:"sell_price"`
	LastTrade string `json:"last_trade"`
	High      string `json:"high"`
	Low       string `json:"low"`
	Avg       string `json:"avg"`
	Vol       string `json:"vol"`
	VolCurr   string `json:"vol_curr"`
	Updated   int64  `json:"updated"`
}
