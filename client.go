package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Exmo struct {
    client *http.Client
    url    string
}
type Currencies map[string]struct{}

func (e *Exmo) GetCurrencies() (Currencies, error) {
    resp, err := e.client.Get(e.url + "/currency")
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("server returned status %d", resp.StatusCode)
    }

    var currencies []string
    if err := json.NewDecoder(resp.Body).Decode(&currencies); err != nil {
        return nil, err
    }

    result := make(Currencies)
    for _, c := range currencies {
        result[c] = struct{}{}
    }
    return result, nil
}


func NewExmo(opts ...func(exmo *Exmo)) Exchanger {
    exmo := &Exmo{
        client: &http.Client{},
        url:    "https://api.exmo.com/v1",
    }
    for _, opt := range opts {
        opt(exmo)
    }
    return exmo
}
func (e *Exmo) GetOrderBook(limit int, pairs ...string) (OrderBook, error) {
    if len(pairs) == 0 {
        return nil, errors.New("at least one pair is required")
    }

    resp, err := e.client.Get(fmt.Sprintf("%s/order_book?limit=%d&pair=%s", e.url, limit, strings.Join(pairs, ",")))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("server returned status %d", resp.StatusCode)
    }

    var orderBook OrderBook
    if err := json.NewDecoder(resp.Body).Decode(&orderBook); err != nil {
        return nil, err
    }

    // Проверяем, что получили данные хотя бы для одной пары
    if len(orderBook) == 0 {
        return nil, errors.New("empty order book response")
    }

    return orderBook, nil
}
func (e *Exmo) GetTicker() (Ticker, error) {
	resp, err := e.client.Get(e.url + "/ticker")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("server returned non-200 status")
	}

	var ticker Ticker
	err = json.NewDecoder(resp.Body).Decode(&ticker)
	return ticker, err
}

func (e *Exmo) GetTrades(pairs ...string) (Trades, error) {
    if len(pairs) == 0 {
        return nil, errors.New("at least one pair is required")
    }

    resp, err := e.client.Get(e.url + "/trades?pair=" + strings.Join(pairs, ","))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("server returned status %d", resp.StatusCode)
    }

    var trades Trades
    if err := json.NewDecoder(resp.Body).Decode(&trades); err != nil {
        return nil, err
    }
    return trades, nil
}



func (e *Exmo) GetClosePrice(pair string, resolution int, start, end time.Time) ([]float64, error) {
	resp, err := e.client.Get(e.url + "/candles_history?pair=" + pair)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("server returned non-200 status")
	}

	var data struct {
		Candles [][]float64 `json:"candles"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	prices := make([]float64, len(data.Candles))
	for i, candle := range data.Candles {
		if len(candle) >= 4 {
			prices[i] = candle[4] // close price
		}
	}
	return prices, nil
}
func (e *Exmo) GetCandlesHistory(pair string, period int, start, end time.Time) (CandlesHistory, error) {
    // Заглушка, чтобы удовлетворить интерфейс
    return CandlesHistory{}, nil
}

type CandlesHistory struct{}

type Exchanger interface {
    GetTicker() (Ticker, error)
    GetTrades(pairs ...string) (Trades, error)
    GetOrderBook(limit int, pairs ...string) (OrderBook, error)
    GetCurrencies() (Currencies, error)
    GetCandlesHistory(pair string, period int, start, end time.Time) (CandlesHistory, error)
    GetClosePrice(pair string, resolution int, start, end time.Time) ([]float64, error)
}

