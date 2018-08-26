package main

import "testing"

const (
	testURL = "https://vine.co/v/eMbVnBFmlbU"
	testID  = "eMbVnBFmlbU"
)

func Test_getMP4URL(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test", args{url: testURL}, "test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getMP4URL(tt.args.url); got != tt.want {
				t.Errorf("getMP4URL() = %v, want %v", got, tt.want)
			}
		})
	}
}
