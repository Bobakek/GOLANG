package main

import (
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewIndicator(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockExchanger := NewMockExchanger(ctrl)

	t.Run("default options", func(t *testing.T) {
		ind := NewIndicator(mockExchanger).(*Indicator)
		
		assert.NotNil(t, ind.calculateSMA)
		assert.NotNil(t, ind.calculateEMA)
		assert.Equal(t, mockExchanger, ind.exchange)
	})

	t.Run("with custom options", func(t *testing.T) {
		customSMA := func(data []float64, period int) []float64 { return []float64{1.0} }
		customEMA := func(data []float64, period int) []float64 { return []float64{2.0} }
		
		ind := NewIndicator(
			mockExchanger,
			WithCalculateSMA(customSMA),
			WithCalculateEMA(customEMA),
		).(*Indicator)
		
		assert.Equal(t, []float64{1.0}, ind.calculateSMA(nil, 0))
		assert.Equal(t, []float64{2.0}, ind.calculateEMA(nil, 0))
		assert.Equal(t, mockExchanger, ind.exchange)
	})
}

func TestSMA(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockExchanger := NewMockExchanger(ctrl)
	ind := NewIndicator(mockExchanger)

	now := time.Now()
	from := now.AddDate(0, 0, -2)
	to := now

	testCases := []struct {
		name         string
		pair        string
		resolution  int
		period      int
		from        time.Time
		to          time.Time
		mockData    []float64
		mockError   error
		expectedLen int
		expectedErr bool
	}{
		{
			name:        "success with enough data",
			pair:       "BTC_USD",
			resolution: 30,
			period:     5,
			from:       from,
			to:         to,
			mockData:   []float64{100, 101, 102, 103, 104, 105},
			expectedLen: 2,
		},
		{
			name:        "success with exact period data",
			pair:       "BTC_USD",
			resolution: 30,
			period:     5,
			from:       from,
			to:         to,
			mockData:   []float64{100, 101, 102, 103, 104},
			expectedLen: 1,
		},
		{
			name:        "not enough data",
			pair:       "BTC_USD",
			resolution: 30,
			period:     5,
			from:       from,
			to:         to,
			mockData:   []float64{100, 101, 102},
			expectedLen: 0,
		},
		{
			name:        "exchange error",
			pair:       "BTC_USD",
			resolution: 30,
			period:     5,
			from:       from,
			to:         to,
			mockError:  errors.New("exchange error"),
			expectedErr: true,
		},
		{
			name:        "empty data",
			pair:       "BTC_USD",
			resolution: 30,
			period:     5,
			from:       from,
			to:         to,
			mockData:   []float64{},
			expectedLen: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockExchanger.EXPECT().
				GetClosePrice(tc.pair, tc.resolution, tc.from, tc.to).
				Return(tc.mockData, tc.mockError)

			result, err := ind.SMA(tc.pair, tc.resolution, tc.period, tc.from, tc.to)
			
			if tc.expectedErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedLen, len(result))
			}
		})
	}
}

func TestEMA(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockExchanger := NewMockExchanger(ctrl)
	ind := NewIndicator(mockExchanger)

	now := time.Now()
	from := now.AddDate(0, 0, -2)
	to := now

	testCases := []struct {
		name         string
		pair        string
		resolution  int
		period      int
		from        time.Time
		to          time.Time
		mockData    []float64
		mockError   error
		expectedErr bool
	}{
		{
			name:        "success with enough data",
			pair:       "BTC_USD",
			resolution: 30,
			period:     5,
			from:       from,
			to:         to,
			mockData:   []float64{100, 101, 102, 103, 104, 105},
		},
		{
			name:        "success with exact period data",
			pair:       "BTC_USD",
			resolution: 30,
			period:     5,
			from:       from,
			to:         to,
			mockData:   []float64{100, 101, 102, 103, 104},
		},
		{
			name:        "not enough data",
			pair:       "BTC_USD",
			resolution: 30,
			period:     5,
			from:       from,
			to:         to,
			mockData:   []float64{100, 101, 102},
		},
		{
			name:        "exchange error",
			pair:       "BTC_USD",
			resolution: 30,
			period:     5,
			from:       from,
			to:         to,
			mockError:  errors.New("exchange error"),
			expectedErr: true,
		},
		{
			name:        "empty data",
			pair:       "BTC_USD",
			resolution: 30,
			period:     5,
			from:       from,
			to:         to,
			mockData:   []float64{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockExchanger.EXPECT().
				GetClosePrice(tc.pair, tc.resolution, tc.from, tc.to).
				Return(tc.mockData, tc.mockError)

			result, err := ind.EMA(tc.pair, tc.resolution, tc.period, tc.from, tc.to)
			
			if tc.expectedErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				if len(tc.mockData) >= tc.period {
					// Для EMA ожидаем результат той же длины, что и входные данные
					assert.Equal(t, len(tc.mockData), len(result))
				} else {
					assert.Empty(t, result)
				}
			}
		})
	}
}

func TestCalculateSMA(t *testing.T) {
	tests := []struct {
		name     string
		data     []float64
		period   int
		expected []float64
	}{
		{
			name:     "basic test",
			data:     []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			period:   3,
			expected: []float64{2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			name:     "empty data",
			data:     []float64{},
			period:   3,
			expected: []float64{},
		},
		{
			name:     "period larger than data",
			data:     []float64{1, 2},
			period:   3,
			expected: []float64{},
		},
		{
			name:     "period equals data length",
			data:     []float64{1, 2, 3},
			period:   3,
			expected: []float64{2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateSMA(tt.data, tt.period)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCalculateEMA(t *testing.T) {
	tests := []struct {
		name     string
		data     []float64
		period   int
		expected []float64
	}{
		{
			name:     "basic test",
			data:     []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			period:   3,
			// Ожидаемые значения для EMA с периодом 3
			expected: []float64{1, 1.5, 2.25, 3.125, 4.0625, 5.03125, 6.015625, 7.0078125, 8.00390625, 9.001953125},
		},
		{
			name:     "empty data",
			data:     []float64{},
			period:   3,
			expected: []float64{},
		},
		{
			name:     "period larger than data",
			data:     []float64{1, 2},
			period:   3,
			expected: []float64{},
		},
		{
			name:     "period equals data length",
			data:     []float64{1, 2, 3},
			period:   3,
			expected: []float64{1, 1.5, 2.25},
		},
		{
			name:     "invalid period (zero)",
			data:     []float64{1, 2, 3},
			period:   0,
			expected: []float64{},
		},
		{
			name:     "invalid period (negative)",
			data:     []float64{1, 2, 3},
			period:   -1,
			expected: []float64{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateEMA(tt.data, tt.period)
			assert.Equal(t, len(tt.expected), len(result), "length mismatch")
			for i := range tt.expected {
				assert.InDelta(t, tt.expected[i], result[i], 0.0001, 
					"at index %d, expected %f, got %f", i, tt.expected[i], result[i])
			}
		})
	}
}