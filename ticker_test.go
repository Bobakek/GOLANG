package main

import (
	"reflect"
	"testing"
)

func TestUnmarshalTicker(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    Ticker
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalTicker(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalTicker() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnmarshalTicker() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTicker_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		r       *Ticker
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.Marshal()
			if (err != nil) != tt.wantErr {
				t.Errorf("Ticker.Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ticker.Marshal() = %v, want %v", got, tt.want)
			}
		})
	}
}
