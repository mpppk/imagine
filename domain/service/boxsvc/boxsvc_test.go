package boxsvc_test

import (
	"reflect"
	"testing"

	"github.com/mpppk/imagine/domain/model"
)

func TestMerge(t *testing.T) {
	type args struct {
		baseBoxes  []*model.BoundingBox
		otherBoxes []*model.BoundingBox
	}
	tests := []struct {
		name         string
		args         args
		wantNewBoxes []*model.BoundingBox
	}{
		{
			name: "append boxes",
			args: args{
				baseBoxes:  []*model.BoundingBox{{TagID: 1}, {TagID: 2}},
				otherBoxes: []*model.BoundingBox{{TagID: 1, X: 1}, {TagID: 2, X: 1}},
			},
			wantNewBoxes: []*model.BoundingBox{
				{TagID: 1}, {TagID: 2}, {TagID: 1, X: 1}, {TagID: 2, X: 1},
			},
		},
		{
			name: "skip if same box",
			args: args{
				baseBoxes:  []*model.BoundingBox{{TagID: 1}, {TagID: 2}},
				otherBoxes: []*model.BoundingBox{{TagID: 1}, {TagID: 2, X: 1}},
			},
			wantNewBoxes: []*model.BoundingBox{
				{TagID: 1}, {TagID: 2}, {TagID: 2, X: 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotNewBoxes := (model.BoundingBoxes)(tt.args.baseBoxes).Merge(tt.args.otherBoxes); !reflect.DeepEqual(gotNewBoxes, tt.wantNewBoxes) {
				t.Errorf("Merge() = %v, want %v", gotNewBoxes, tt.wantNewBoxes)
			}
		})
	}
}
