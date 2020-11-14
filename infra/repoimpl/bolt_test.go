package repoimpl

import (
	"testing"
)

func Test_itob(t *testing.T) {
	tests := []struct {
		name string
		v    uint64
	}{
		{v: 0},
		{v: 1},
		{v: 8},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := btoi(itob(tt.v)); result != tt.v {
				t.Errorf("provided: %v, got: %v", tt.v, result)
			}
		})
	}
}
