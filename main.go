package main

import (
	"fmt"
	"time"
)

// Добавляем глобальную переменную для возможности подмены в тестах
var globalExchanger Exchanger = NewExmo()

func main() {
	// Используем globalExchanger вместо создания нового экземпляра
	indicator := NewIndicator(globalExchanger,
		WithCalculateSMA(calculateSMA),
		WithCalculateEMA(calculateEMA),
	)

	// Остальной код остается без изменений
	pair := "BTC_USD"
	resolution := 30
	period := 5
	from := time.Now().AddDate(0, 0, -2)
	to := time.Now()

	sma, err := indicator.SMA(pair, resolution, period, from, to)
	if err != nil {
		fmt.Println("Ошибка при расчете SMA:", err)
		return
	}
	fmt.Println("SMA:", sma)

	ema, err := indicator.EMA(pair, resolution, period, from, to)
	if err != nil {
		fmt.Println("Ошибка при расчете EMA:", err)
		return
	}
	fmt.Println("EMA:", ema)
}