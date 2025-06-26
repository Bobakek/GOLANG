package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExmo_GetTicker(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/ticker", r.URL.Path)
		w.Write([]byte(`{"BTC_USD":{"buy_price":"50000","sell_price":"50100"}}`))
	}))
	defer ts.Close()

	client := NewExmo(func(e *Exmo) { e.url = ts.URL })
	ticker, err := client.GetTicker()

	assert.NoError(t, err)
	assert.Equal(t, "50000", ticker["BTC_USD"].BuyPrice)
	assert.Equal(t, "50100", ticker["BTC_USD"].SellPrice)
}

func TestExmo_GetTicker_Error(t *testing.T) {
	t.Run("server error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer ts.Close()

		client := NewExmo(func(e *Exmo) { e.url = ts.URL })
		_, err := client.GetTicker()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "non-200 status")
	})

	t.Run("invalid json", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`invalid`))
		}))
		defer ts.Close()

		client := NewExmo(func(e *Exmo) { e.url = ts.URL })
		_, err := client.GetTicker()
		assert.Error(t, err)
	})
}

func TestExmo_GetClosePrice(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/candles_history", r.URL.Path)
		w.Write([]byte(`{"candles":[[1640995200000,50000,51000,49000,50500,10]]}`))
	}))
	defer ts.Close()

	client := NewExmo(func(e *Exmo) { e.url = ts.URL })
	prices, err := client.GetClosePrice("BTC_USD", 60, time.Now(), time.Now())
	
	assert.NoError(t, err)
	assert.Equal(t, []float64{50500}, prices)
}

func TestExmo_GetClosePrice_Error(t *testing.T) {
	t.Run("server error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer ts.Close()

		client := NewExmo(func(e *Exmo) { e.url = ts.URL })
		_, err := client.GetClosePrice("BTC_USD", 60, time.Now(), time.Now())
		assert.Error(t, err)
	})

	t.Run("invalid json", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`invalid`))
		}))
		defer ts.Close()

		client := NewExmo(func(e *Exmo) { e.url = ts.URL })
		_, err := client.GetClosePrice("BTC_USD", 60, time.Now(), time.Now())
		assert.Error(t, err)
	})
}

func TestExmo_GetTrades(t *testing.T) {
    t.Run("success", func(t *testing.T) {
        ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            assert.Equal(t, "/trades", r.URL.Path)
            assert.Equal(t, "pair=BTC_USD", r.URL.RawQuery)
            w.Write([]byte(`{"BTC_USD":[{"trade_id":1,"price":"50000"}]}`))
        }))
        defer ts.Close()

        client := NewExmo(func(e *Exmo) { e.url = ts.URL })
        trades, err := client.GetTrades("BTC_USD")
        
        assert.NoError(t, err)
        assert.Equal(t, int64(1), trades["BTC_USD"][0].TradeID)
    })

    t.Run("empty pairs", func(t *testing.T) {
        client := NewExmo()
        _, err := client.GetTrades()
        assert.Error(t, err)
        assert.Contains(t, err.Error(), "at least one pair is required")
    })

    t.Run("server error", func(t *testing.T) {
        ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(http.StatusInternalServerError)
        }))
        defer ts.Close()

        client := NewExmo(func(e *Exmo) { e.url = ts.URL })
        _, err := client.GetTrades("BTC_USD")
        assert.Error(t, err)
    })

    t.Run("invalid json", func(t *testing.T) {
        ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Write([]byte(`invalid`))
        }))
        defer ts.Close()

        client := NewExmo(func(e *Exmo) { e.url = ts.URL })
        _, err := client.GetTrades("BTC_USD")
        assert.Error(t, err)
    })
}

func TestExmo_GetCurrencies(t *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        assert.Equal(t, "/currency", r.URL.Path)
        w.Write([]byte(`["BTC","USD"]`))
    }))
    defer ts.Close()

    client := NewExmo(func(e *Exmo) { e.url = ts.URL })
    currencies, err := client.GetCurrencies()
    
    assert.NoError(t, err)
    _, btcExists := currencies["BTC"]
    assert.True(t, btcExists)
}

func TestExmo_GetCandlesHistory(t *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        assert.Equal(t, "/candles_history", r.URL.Path)
        w.Write([]byte(`{"candles":[[1640995200000,50000,51000,49000,50500,10]]}`))
    }))
    defer ts.Close()

    client := NewExmo(func(e *Exmo) { e.url = ts.URL })
    _, err := client.GetCandlesHistory("BTC_USD", 60, time.Now(), time.Now())
    assert.NoError(t, err)
}

func TestExmo_GetTicker_EmptyResponse(t *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte(`{}`)) // Пустой ответ
    }))
    defer ts.Close()

    client := NewExmo(func(e *Exmo) { e.url = ts.URL })
    ticker, err := client.GetTicker()
    
    assert.NoError(t, err)
    assert.Empty(t, ticker)
}

func TestExmo_GetClosePrice_EmptyCandles(t *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte(`{"candles":[]}`)) // Пустые свечи
    }))
    defer ts.Close()

    client := NewExmo(func(e *Exmo) { e.url = ts.URL })
    prices, err := client.GetClosePrice("BTC_USD", 60, time.Now(), time.Now())
    
    assert.NoError(t, err)
    assert.Empty(t, prices)
}

func TestExmo_GetOrderBook(t *testing.T) {
    t.Run("success", func(t *testing.T) {
        ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            assert.Equal(t, "/order_book", r.URL.Path)
            assert.Equal(t, "limit=10&pair=BTC_USD", r.URL.RawQuery)
            w.Write([]byte(`{"BTC_USD":{"ask":[["50000","1"]],"bid":[["49000","1"]]}}`))
        }))
        defer ts.Close()

        client := NewExmo(func(e *Exmo) { e.url = ts.URL })
        book, err := client.GetOrderBook(10, "BTC_USD")
        
        assert.NoError(t, err)
        assert.Equal(t, "50000", book["BTC_USD"].Ask[0][0])
        assert.Equal(t, "49000", book["BTC_USD"].Bid[0][0])
    })

    t.Run("empty pairs", func(t *testing.T) {
        client := NewExmo()
        _, err := client.GetOrderBook(10)
        assert.Error(t, err)
        assert.Contains(t, err.Error(), "at least one pair is required")
    })

    t.Run("server error", func(t *testing.T) {
        ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(http.StatusInternalServerError)
        }))
        defer ts.Close()

        client := NewExmo(func(e *Exmo) { e.url = ts.URL })
        _, err := client.GetOrderBook(10, "BTC_USD")
        assert.Error(t, err)
    })

    t.Run("invalid json", func(t *testing.T) {
        ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Write([]byte(`invalid`))
        }))
        defer ts.Close()

        client := NewExmo(func(e *Exmo) { e.url = ts.URL })
        _, err := client.GetOrderBook(10, "BTC_USD")
        assert.Error(t, err)
    })

    t.Run("empty response", func(t *testing.T) {
        ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Write([]byte(`{}`))
        }))
        defer ts.Close()

        client := NewExmo(func(e *Exmo) { e.url = ts.URL })
        _, err := client.GetOrderBook(10, "BTC_USD")
        assert.Error(t, err)
        assert.Contains(t, err.Error(), "empty order book response")
    })
}

