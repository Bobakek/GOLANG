package main

import "encoding/json"

func UnmarshalTrades(data []byte) (Trades, error) {
	var r Trades
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Trades) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Trades map[string][]Pair

type Pair struct {
	TradeID  int64  `json:"trade_id"`
	Date     int64  `json:"date"`
	Type     Type   `json:"type"`
	Quantity string `json:"quantity"`
	Price    string `json:"price"`
	Amount   string `json:"amount"`
}

type Type string

const (
	Buy  Type = "buy"
	Sell Type = "sell"
)
