package interactor_test

import (
	"testing"

	"github.com/mpppk/imagine/testutil"

	"github.com/mpppk/imagine/usecase/usecasetest"

	"github.com/mpppk/imagine/domain/model"
)

func TestTag_SetTags(t *testing.T) {
	type args struct {
		ws       model.WSName
		tagNames []string
	}
	tests := []struct {
		name      string
		existTags []*model.Tag
		want      []*model.Tag
		args      args
		wantErr   bool
	}{
		{
			name:      "must be empty",
			existTags: []*model.Tag{{ID: 1, Name: "tag1"}},
			args:      args{ws: "default-workspace", tagNames: []string{}},
			want:      nil,
		},
		{
			name:      "set one tag",
			existTags: []*model.Tag{{ID: 1, Name: "tag1"}},
			args:      args{ws: "default-workspace", tagNames: []string{"new-tag1"}},
			want:      []*model.Tag{{ID: 1, Name: "new-tag1"}},
		},
		{
			name:      "set two tags",
			existTags: []*model.Tag{{ID: 1, Name: "tag1"}},
			args:      args{ws: "default-workspace", tagNames: []string{"new-tag1", "new-tag2"}},
			want:      []*model.Tag{{ID: 1, Name: "new-tag1"}, {ID: 2, Name: "new-tag2"}},
		},
	}
	for _, tt := range tests {
		tt := tt
		usecasetest.RunParallelWithUseCases(t, tt.name, tt.args.ws, func(t *testing.T, ut *usecasetest.UseCases) {
			idList, err := ut.Usecases.Tag.SetTags(tt.args.ws, tt.args.tagNames)
			if (err != nil) != tt.wantErr {
				t.Errorf("PutTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if tt.wantErr {
				return
			}
			var wantIdList []model.TagID
			for _, tag := range tt.want {
				wantIdList = append(wantIdList, tag.ID)
			}

			testutil.Diff(t, wantIdList, idList)

			tags := ut.Tag.List(tt.args.ws)
			testutil.Diff(t, tt.want, tags)
		})
	}
}

func TestTag_SetTagsTwice(t *testing.T) {
	type args struct {
		ws       model.WSName
		tagNames []string
	}
	tests := []struct {
		name      string
		existTags []*model.Tag
		want1     []*model.Tag
		want2     []*model.Tag
		args1     args
		args2     args
		wantErr   bool
	}{
		{
			name:      "Same name tag should have same ID even if once deleted",
			existTags: []*model.Tag{{ID: 1, Name: "tag1"}, {ID: 2, Name: "tag2"}},
			args1:     args{ws: "default-workspace", tagNames: []string{}},
			args2:     args{ws: "default-workspace", tagNames: []string{"tag2"}},
			want1:     nil,
			want2:     []*model.Tag{{ID: 2, Name: "new-tag2"}},
		},
		{
			name:      "Different name tag should have different ID",
			existTags: []*model.Tag{{ID: 1, Name: "tag1"}},
			args1:     args{ws: "default-workspace", tagNames: []string{}},
			args2:     args{ws: "default-workspace", tagNames: []string{"tag2"}},
			want1:     nil,
			want2:     []*model.Tag{{ID: 2, Name: "new-tag2"}},
		},
	}
	testF := func(t *testing.T, ut *usecasetest.UseCases, args args, want []*model.Tag, wantErr bool) {
		idList, err := ut.Usecases.Tag.SetTags(args.ws, args.tagNames)
		if (err != nil) != wantErr {
			t.Errorf("PutTags() error = %v, wantErr %v", err, wantErr)
			return
		} else if wantErr {
			return
		}
		var wantIdList []model.TagID
		for _, tag := range want {
			wantIdList = append(wantIdList, tag.ID)
		}

		testutil.Diff(t, wantIdList, idList)

		tags := ut.Tag.List(args.ws)
		testutil.Diff(t, want, tags)
	}
	for _, tt := range tests {
		tt := tt
		usecasetest.RunParallelWithUseCases(t, tt.name, tt.args1.ws, func(t *testing.T, ut *usecasetest.UseCases) {
			testF(t, ut, tt.args1, tt.want1, tt.wantErr)
			testF(t, ut, tt.args2, tt.want2, tt.wantErr)
			//idList, err := ut.Usecases.Tag.SetTags(tt.args1.ws, tt.args1.tagNames)
			//if (err != nil) != tt.wantErr {
			//	t.Errorf("PutTags() error = %v, wantErr %v", err, tt.wantErr)
			//	return
			//} else if tt.wantErr {
			//	return
			//}
			//var wantIdList []model.TagID
			//for _, tag := range tt.want1 {
			//	wantIdList = append(wantIdList, tag.ID)
			//}
			//
			//testutil.Diff(t, wantIdList, idList)
			//
			//tags := ut.Tag.List(tt.args1.ws)
			//testutil.Diff(t, tt.want1, tags)
		})
	}
}
