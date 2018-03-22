package main

import (
	"testing"
)

func TestEnsureDir(t *testing.T) {
	tests := []struct {
		name         string
		destionation string
	}{
		{
			name:         "empty",
			destionation: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ensureDir(tt.destionation)
		})
	}
}
