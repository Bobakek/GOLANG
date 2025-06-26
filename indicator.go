package main

import (
	"time"

	"github.com/cinar/indicator"
)

type Indicatorer interface {
	SMA(pair string, resolution, period int, from, to time.Time) ([]float64, error)
	EMA(pair string, resolution, period int, from, to time.Time) ([]float64, error)
}

type Indicator struct {
	exchange     Exchanger
	calculateSMA func(data []float64, period int) []float64
	calculateEMA func(data []float64, period int) []float64
}

type IndicatorOption func(*Indicator)

func WithCalculateSMA(f func(data []float64, period int) []float64) IndicatorOption {
	return func(i *Indicator) {
		i.calculateSMA = f
	}
}

func WithCalculateEMA(f func(data []float64, period int) []float64) IndicatorOption {
	return func(i *Indicator) {
		i.calculateEMA = f
	}
}

func NewIndicator(exchange Exchanger, opts ...IndicatorOption) Indicatorer {
	ind := &Indicator{
		exchange: exchange,
	}

	for _, opt := range opts {
		opt(ind)
	}

	// Установка значений по умолчанию
	if ind.calculateSMA == nil {
		ind.calculateSMA = calculateSMA
	}

	if ind.calculateEMA == nil {
		ind.calculateEMA = calculateEMA
	}

	return ind
}

func (i *Indicator) SMA(pair string, resolution, period int, from, to time.Time) ([]float64, error) {
	data, err := i.exchange.GetClosePrice(pair, resolution, from, to)
	if err != nil {
		return nil, err
	}

	return i.calculateSMA(data, period), nil
}

func (i *Indicator) EMA(pair string, resolution, period int, from, to time.Time) ([]float64, error) {
	data, err := i.exchange.GetClosePrice(pair, resolution, from, to)
	if err != nil {
		return nil, err
	}

	return i.calculateEMA(data, period), nil
}

func calculateSMA(data []float64, period int) []float64 {
    if len(data) < period || period <= 0 {
        return []float64{}
    }
    return indicator.Sma(period, data)[period-1:]
}

func calculateEMA(data []float64, period int) []float64 {
    if len(data) < period || period <= 0 {
        return []float64{}
    }
    // EMA возвращает результат той же длины, что и входные данные
    return indicator.Ema(period, data)
}