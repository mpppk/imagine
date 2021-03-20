package queryimpl_test

import (
	"sort"
	"testing"

	"github.com/mpppk/imagine/testutil"
	"github.com/mpppk/imagine/usecase/usecasetest"

	"github.com/mpppk/imagine/domain/model"
)

func TestBBoltTag_ListByQueries(t *testing.T) {
	type args struct {
		ws      model.WSName
		queries []*model.Query
	}
	tests := []struct {
		name          string
		ws            model.WSName
		existTagNames []string
		args          args
		wantTags      []*model.Tag
		wantErr       bool
	}{
		{
			name:          "list all tags if queries is nil",
			ws:            "default-workspace",
			existTagNames: []string{"tag1", "tag2", "tag3"},
			args: args{
				ws:      "default-workspace",
				queries: nil,
			},
			wantTags: []*model.Tag{
				testutil.MustNewTag(1, "tag1", 0),
				testutil.MustNewTag(2, "tag2", 1),
				testutil.MustNewTag(3, "tag3", 2),
			},
			wantErr: false,
		},
		{
			name:          "list only query matched tags",
			ws:            "default-workspace",
			existTagNames: []string{"tag1", "tag2", "tag3"},
			args: args{
				ws: "default-workspace",
				queries: []*model.Query{
					{Op: "equals", Value: "tag1"},
				},
			},
			wantTags: []*model.Tag{
				testutil.MustNewTag(1, "tag1", 0),
			},
			wantErr: false,
		},
		{
			name:          "any tag does not match if unsupported query(path-equals) is specified",
			ws:            "default-workspace",
			existTagNames: []string{"tag1", "tag2", "tag3"},
			args: args{
				ws: "default-workspace",
				queries: []*model.Query{
					// path-equals query can't use to tag search
					{Op: "path-equals", Value: "path1"},
				},
			},
			wantTags: nil,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		tt := tt
		usecasetest.RunParallelWithUseCases(t, tt.name, tt.args.ws, func(t *testing.T, ut *usecasetest.UseCases) {
			ut.Tag.SetTagByNames(tt.args.ws, tt.existTagNames)
			gotTags, err := ut.Usecases.Client.Tag.ListByQueries(tt.args.ws, tt.args.queries)
			if (err != nil) != tt.wantErr {
				t.Errorf("PutTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if tt.wantErr {
				return
			}
			sort.Slice(gotTags, func(i, j int) bool {
				return gotTags[i].ID < gotTags[j].ID
			})
			testutil.Diff(t, tt.wantTags, gotTags)
		})
	}
}
