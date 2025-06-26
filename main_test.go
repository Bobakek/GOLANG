package main

import (
	"io"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMainOutput(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockExchanger := NewMockExchanger(ctrl)
	
	// Используем gomock.Any() для временных параметров
	mockExchanger.EXPECT().GetClosePrice(
		"BTC_USD", 
		30, 
		gomock.Any(), // Любое время начала
		gomock.Any(), // Любое время окончания
	).Return([]float64{100, 101, 102, 103, 104}, nil).Times(2)

	// Перехватываем stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Подменяем globalExchanger
	original := globalExchanger
	globalExchanger = mockExchanger
	defer func() { 
		globalExchanger = original 
		w.Close()
		os.Stdout = old
	}()

	main()

	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = old

	output := string(out)
	
	// Проверяем общую структуру вывода
	assert.Contains(t, output, "SMA:")
	assert.Contains(t, output, "EMA:")
	
	// Проверяем что вывод содержит списки чисел
	assert.Regexp(t, `SMA: \[[\d\., ]+\]`, output)
	assert.Regexp(t, `EMA: \[[\d\., ]+\]`, output)
}