package main

import (
	
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderBookMarshaling(t *testing.T) {
	t.Run("unmarshal valid data", func(t *testing.T) {
		data := []byte(`{"BTC_USD":{"ask":[["50000","1"]],"bid":[["49000","1"]]}}`)
		ob, err := UnmarshalOrderBook(data)
		
		assert.NoError(t, err)
		assert.Equal(t, 1, len(ob["BTC_USD"].Ask))
		assert.Equal(t, "50000", ob["BTC_USD"].Ask[0][0])
	})

	t.Run("marshal empty data", func(t *testing.T) {
		ob := make(OrderBook)
		data, err := ob.Marshal()
		
		assert.NoError(t, err)
		assert.Equal(t, `{}`, string(data))
	})

	t.Run("unmarshal invalid data", func(t *testing.T) {
		_, err := UnmarshalOrderBook([]byte(`invalid`))
		assert.Error(t, err)
	})
}