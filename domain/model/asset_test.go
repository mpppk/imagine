package model

import (
	"reflect"
	"testing"
)

func TestReplaceBoundingBoxByID(t *testing.T) {
	type args struct {
		boxes      []*BoundingBox
		replaceBox *BoundingBox
	}
	tests := []struct {
		name         string
		args         args
		wantNewBoxes []*BoundingBox
	}{
		{
			args: args{
				boxes: []*BoundingBox{
					{ID: 0},
					{ID: 1},
					{ID: 2},
				},
				replaceBox: &BoundingBox{ID: 1, X: 1},
			},
			wantNewBoxes: []*BoundingBox{
				{ID: 0}, {ID: 1, X: 1}, {ID: 2},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotNewBoxes := ReplaceBoundingBoxByID(tt.args.boxes, tt.args.replaceBox); !reflect.DeepEqual(gotNewBoxes, tt.wantNewBoxes) {
				t.Errorf("ReplaceBoundingBoxByID() = %v, want %v", gotNewBoxes, tt.wantNewBoxes)
			}
		})
	}
}
