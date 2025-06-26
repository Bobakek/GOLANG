package main

import (
	
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTradesMarshaling(t *testing.T) {
	t.Run("unmarshal valid trade", func(t *testing.T) {
		data := []byte(`{"BTC_USD":[{"trade_id":1,"type":"buy","price":"50000"}]}`)
		trades, err := UnmarshalTrades(data)
		
		assert.NoError(t, err)
		assert.Equal(t, Buy, trades["BTC_USD"][0].Type)
		assert.Equal(t, "50000", trades["BTC_USD"][0].Price)
	})

	t.Run("marshal trade", func(t *testing.T) {
		trades := Trades{
			"BTC_USD": []Pair{{TradeID: 1, Type: Buy}},
		}
		data, err := trades.Marshal()
		
		assert.NoError(t, err)
		assert.Contains(t, string(data), `"trade_id":1`)
	})
}