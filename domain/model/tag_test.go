package model_test

import (
	"testing"

	"github.com/mpppk/imagine/domain/model"

	"github.com/mpppk/imagine/testutil"
)

func TestNewTagSet(t *testing.T) {
	type args struct {
		tags []*model.Tag
	}
	tests := []struct {
		name      string
		args      args
		wantM     map[model.TagID]*model.Tag
		wantNameM map[string]*model.Tag
	}{
		{
			name:      "return empty TagSet from empty tags",
			args:      args{tags: []*model.Tag{}},
			wantM:     map[model.TagID]*model.Tag{},
			wantNameM: map[string]*model.Tag{},
		},
		{
			name: "return TagSet from one tag",
			args: args{tags: []*model.Tag{
				testutil.MustNewTag(1, "tag1", 0),
			}},
			wantM: map[model.TagID]*model.Tag{
				1: testutil.MustNewTag(1, "tag1", 0),
			},
			wantNameM: map[string]*model.Tag{
				"tag1": testutil.MustNewTag(1, "tag1", 0),
			},
		},
		{
			name: "return TagSet from two tags",
			args: args{tags: []*model.Tag{
				testutil.MustNewTag(1, "tag1", 0),
				testutil.MustNewTag(2, "tag2", 1),
			}},
			wantM: map[model.TagID]*model.Tag{
				1: testutil.MustNewTag(1, "tag1", 0),
				2: testutil.MustNewTag(2, "tag2", 1),
			},
			wantNameM: map[string]*model.Tag{
				"tag1": testutil.MustNewTag(1, "tag1", 0),
				"tag2": testutil.MustNewTag(2, "tag2", 1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tagSet := model.NewTagSet(tt.args.tags)
			m, nameM := tagSet.ToMap()
			testutil.Diff(t, tt.wantM, m)
			testutil.Diff(t, tt.wantNameM, nameM)
		})
	}
}

func TestTagSet_Get(t1 *testing.T) {
	type args struct {
		id model.TagID
	}
	tests := []struct {
		name   string
		tagSet *model.TagSet
		args   args
		want   *model.Tag
		want1  bool
	}{
		{
			name:   "get from empty TagSet",
			tagSet: model.NewTagSet([]*model.Tag{}),
			args:   args{id: 1},
			want:   nil,
			want1:  false,
		},
		{
			name: "get exist tag from one tag",
			tagSet: model.NewTagSet([]*model.Tag{
				testutil.MustNewTag(1, "tag1", 0),
			}),
			args:  args{id: 1},
			want:  testutil.MustNewTag(1, "tag1", 0),
			want1: true,
		},
		{
			name: "get non exist tag from one tag",
			tagSet: model.NewTagSet([]*model.Tag{
				testutil.MustNewTag(1, "tag1", 0),
			}),
			args:  args{id: 2},
			want:  nil,
			want1: false,
		},
		{
			name: "get exist tag from two tags",
			tagSet: model.NewTagSet([]*model.Tag{
				testutil.MustNewTag(1, "tag1", 0),
				testutil.MustNewTag(2, "tag2", 1),
			}),
			args:  args{id: 2},
			want:  testutil.MustNewTag(2, "tag2", 1),
			want1: true,
		},
		{
			name: "get exist tag from two tags",
			tagSet: model.NewTagSet([]*model.Tag{
				testutil.MustNewTag(1, "tag1", 0),
				testutil.MustNewTag(2, "tag2", 1),
			}),
			args:  args{id: 3},
			want:  nil,
			want1: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			got, got1 := tt.tagSet.Get(tt.args.id)
			testutil.Diff(t1, tt.want, got)
			testutil.Diff(t1, tt.want1, got1)
		})
	}
}

func TestTagSet_GetByName(t1 *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		tagSet *model.TagSet
		args   args
		want   *model.Tag
		want1  bool
	}{
		{
			name:   "get from empty TagSet",
			tagSet: model.NewTagSet([]*model.Tag{}),
			args:   args{name: "tag1"},
			want:   nil,
			want1:  false,
		},
		{
			name: "get exist tag from one tag",
			tagSet: model.NewTagSet([]*model.Tag{
				testutil.MustNewTag(1, "tag1", 0),
			}),
			args:  args{name: "tag1"},
			want:  testutil.MustNewTag(1, "tag1", 0),
			want1: true,
		},
		{
			name: "get non exist tag from one tag",
			tagSet: model.NewTagSet([]*model.Tag{
				testutil.MustNewTag(1, "tag1", 0),
			}),
			args:  args{name: "tag2"},
			want:  nil,
			want1: false,
		},
		{
			name: "get exist tag from two tags",
			tagSet: model.NewTagSet([]*model.Tag{
				testutil.MustNewTag(1, "tag1", 0),
				testutil.MustNewTag(2, "tag2", 1),
			}),
			args:  args{name: "tag2"},
			want:  testutil.MustNewTag(2, "tag2", 1),
			want1: true,
		},
		{
			name: "get exist tag from two tags",
			tagSet: model.NewTagSet([]*model.Tag{
				testutil.MustNewTag(1, "tag1", 0),
				testutil.MustNewTag(2, "tag2", 1),
			}),
			args:  args{name: "tag3"},
			want:  nil,
			want1: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			got, got1 := tt.tagSet.GetByName(tt.args.name)
			testutil.Diff(t1, tt.want, got)
			testutil.Diff(t1, tt.want1, got1)
		})
	}
}

func TestTagSet_Set(t1 *testing.T) {
	type args struct {
		tag *model.Tag
	}
	tests := []struct {
		name      string
		tagSet    *model.TagSet
		args      args
		want      bool
		wantM     map[model.TagID]*model.Tag
		wantNameM map[string]*model.Tag
	}{
		{
			name: "add tag to TagSet",
			tagSet: model.NewTagSet([]*model.Tag{
				testutil.MustNewTag(1, "tag1", 0),
				testutil.MustNewTag(2, "tag2", 1),
			}),
			args: args{tag: testutil.MustNewTag(3, "tag3", 2)},
			want: true,
			wantM: map[model.TagID]*model.Tag{
				1: testutil.MustNewTag(1, "tag1", 0),
				2: testutil.MustNewTag(2, "tag2", 1),
				3: testutil.MustNewTag(3, "tag3", 2),
			},
			wantNameM: map[string]*model.Tag{
				"tag1": testutil.MustNewTag(1, "tag1", 0),
				"tag2": testutil.MustNewTag(2, "tag2", 1),
				"tag3": testutil.MustNewTag(3, "tag3", 2),
			},
		},
		{
			name: "update tag",
			tagSet: model.NewTagSet([]*model.Tag{
				testutil.MustNewTag(1, "tag1", 0),
				testutil.MustNewTag(2, "tag2", 1),
			}),
			args: args{tag: testutil.MustNewTag(2, "updated-tag2", 0)},
			want: true,
			wantM: map[model.TagID]*model.Tag{
				1: testutil.MustNewTag(1, "tag1", 0),
				2: testutil.MustNewTag(2, "updated-tag2", 1),
			},
			wantNameM: map[string]*model.Tag{
				"tag1": testutil.MustNewTag(1, "tag1", 0),
				"tag2": testutil.MustNewTag(2, "updated-tag2", 1),
			},
		},
		{
			name: "do nothing for same tag",
			tagSet: model.NewTagSet([]*model.Tag{
				testutil.MustNewTag(1, "tag1", 0),
				testutil.MustNewTag(2, "tag2", 1),
			}),
			args: args{tag: testutil.MustNewTag(2, "tag2", 1)},
			want: true,
			wantM: map[model.TagID]*model.Tag{
				1: testutil.MustNewTag(1, "tag1", 0),
				2: testutil.MustNewTag(2, "tag2", 1),
			},
			wantNameM: map[string]*model.Tag{
				"tag1": testutil.MustNewTag(1, "tag1", 0),
				"tag2": testutil.MustNewTag(2, "tag2", 1),
			},
		},
		{
			name: "do nothing to tag which have same name but different ID, then return false",
			tagSet: model.NewTagSet([]*model.Tag{
				testutil.MustNewTag(1, "tag1", 0),
				testutil.MustNewTag(2, "tag2", 1),
			}),
			args: args{tag: testutil.MustNewTag(3, "tag2", 1)},
			want: false,
			wantM: map[model.TagID]*model.Tag{
				1: testutil.MustNewTag(1, "tag1", 0),
				2: testutil.MustNewTag(2, "tag2", 1),
			},
			wantNameM: map[string]*model.Tag{
				"tag1": testutil.MustNewTag(1, "tag1", 0),
				"tag2": testutil.MustNewTag(2, "tag2", 1),
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			got := tt.tagSet.Set(tt.args.tag)
			testutil.Diff(t1, tt.want, got)
		})
	}
}

func TestTagSet_SubSetBy(t1 *testing.T) {
	type args struct {
		f func(tag *model.Tag) bool
	}
	tests := []struct {
		name      string
		tagSet    *model.TagSet
		args      args
		wantM     map[model.TagID]*model.Tag
		wantNameM map[string]*model.Tag
	}{
		{
			name:      "return empty subset if TagSet is empty",
			tagSet:    model.NewTagSet([]*model.Tag{}),
			args:      args{f: func(tag *model.Tag) bool { return true }},
			wantM:     map[model.TagID]*model.Tag{},
			wantNameM: map[string]*model.Tag{},
		},
		{
			name: "return empty subset if any tag does not match",
			tagSet: model.NewTagSet([]*model.Tag{
				testutil.MustNewTag(1, "tag1", 0),
				testutil.MustNewTag(2, "tag2", 1),
				testutil.MustNewTag(3, "tag3", 2),
			}),
			args:      args{f: func(tag *model.Tag) bool { return false }},
			wantM:     map[model.TagID]*model.Tag{},
			wantNameM: map[string]*model.Tag{},
		},
		{
			name: "return subset with two tags",
			tagSet: model.NewTagSet([]*model.Tag{
				testutil.MustNewTag(1, "tag1", 0),
				testutil.MustNewTag(2, "tag2", 1),
				testutil.MustNewTag(3, "tag3", 2),
			}),
			args: args{f: func(tag *model.Tag) bool {
				return tag.Name == "tag1" || tag.Name == "tag2"
			}},
			wantM: map[model.TagID]*model.Tag{
				1: testutil.MustNewTag(1, "tag1", 0),
				2: testutil.MustNewTag(2, "tag2", 1),
			},
			wantNameM: map[string]*model.Tag{
				"tag1": testutil.MustNewTag(1, "tag1", 0),
				"tag2": testutil.MustNewTag(2, "tag2", 1),
			},
		},
		{
			name: "return subset with one tag",
			tagSet: model.NewTagSet([]*model.Tag{
				testutil.MustNewTag(1, "tag1", 0),
				testutil.MustNewTag(2, "tag2", 1),
				testutil.MustNewTag(3, "tag3", 2),
			}),
			args: args{f: func(tag *model.Tag) bool {
				return tag.Name == "tag1"
			}},
			wantM: map[model.TagID]*model.Tag{
				1: testutil.MustNewTag(1, "tag1", 0),
			},
			wantNameM: map[string]*model.Tag{
				"tag1": testutil.MustNewTag(1, "tag1", 0),
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := tt.tagSet
			subset := t.SubSetBy(tt.args.f)
			gotM, gotNameM := subset.ToMap()
			testutil.Diff(t1, tt.wantM, gotM)
			testutil.Diff(t1, tt.wantNameM, gotNameM)
		})
	}
}

func TestTagSet_SplitBy(t1 *testing.T) {
	type args struct {
		f func(tag *model.Tag) bool
	}
	tests := []struct {
		name       string
		tagSet     *model.TagSet
		args       args
		wantTrueM  map[model.TagID]*model.Tag
		wantFalseM map[model.TagID]*model.Tag
	}{
		{
			name:       "return empty subset if TagSet is empty",
			tagSet:     model.NewTagSet([]*model.Tag{}),
			args:       args{f: func(tag *model.Tag) bool { return true }},
			wantTrueM:  map[model.TagID]*model.Tag{},
			wantFalseM: map[model.TagID]*model.Tag{},
		},
		{
			name: "all false",
			tagSet: model.NewTagSet([]*model.Tag{
				testutil.MustNewTag(1, "tag1", 0),
				testutil.MustNewTag(2, "tag2", 1),
				testutil.MustNewTag(3, "tag3", 2),
			}),
			args:      args{f: func(tag *model.Tag) bool { return false }},
			wantTrueM: map[model.TagID]*model.Tag{},
			wantFalseM: map[model.TagID]*model.Tag{
				1: testutil.MustNewTag(1, "tag1", 0),
				2: testutil.MustNewTag(2, "tag2", 1),
				3: testutil.MustNewTag(3, "tag3", 2),
			},
		},
		{
			name: "two true",
			tagSet: model.NewTagSet([]*model.Tag{
				testutil.MustNewTag(1, "tag1", 0),
				testutil.MustNewTag(2, "tag2", 1),
				testutil.MustNewTag(3, "tag3", 2),
			}),
			args: args{f: func(tag *model.Tag) bool {
				return tag.Name == "tag1" || tag.Name == "tag2"
			}},
			wantTrueM: map[model.TagID]*model.Tag{
				1: testutil.MustNewTag(1, "tag1", 0),
				2: testutil.MustNewTag(2, "tag2", 1),
			},
			wantFalseM: map[model.TagID]*model.Tag{
				3: testutil.MustNewTag(3, "tag3", 2),
			},
		},
		{
			name: "one true",
			tagSet: model.NewTagSet([]*model.Tag{
				testutil.MustNewTag(1, "tag1", 0),
				testutil.MustNewTag(2, "tag2", 1),
				testutil.MustNewTag(3, "tag3", 2),
			}),
			args: args{f: func(tag *model.Tag) bool {
				return tag.Name == "tag1"
			}},
			wantTrueM: map[model.TagID]*model.Tag{
				1: testutil.MustNewTag(1, "tag1", 0),
			},
			wantFalseM: map[model.TagID]*model.Tag{
				2: testutil.MustNewTag(2, "tag2", 1),
				3: testutil.MustNewTag(3, "tag3", 2),
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := tt.tagSet
			trueSet, falseSet := t.SplitBy(tt.args.f)
			trueM, _ := trueSet.ToMap()
			falseM, _ := falseSet.ToMap()
			testutil.Diff(t1, tt.wantTrueM, trueM)
			testutil.Diff(t1, tt.wantFalseM, falseM)
		})
	}
}

func TestTagSet_SubSetByID(t1 *testing.T) {
	type args struct {
		idList []model.TagID
	}
	tests := []struct {
		name      string
		tagSet    *model.TagSet
		args      args
		wantM     map[model.TagID]*model.Tag
		wantNameM map[string]*model.Tag
	}{
		{
			name:      "return empty subset if TagSet is empty",
			tagSet:    model.NewTagSet([]*model.Tag{}),
			args:      args{idList: []model.TagID{1}},
			wantM:     map[model.TagID]*model.Tag{},
			wantNameM: map[string]*model.Tag{},
		},
		{
			name: "return empty subset if any tag does not match",
			tagSet: model.NewTagSet([]*model.Tag{
				testutil.MustNewTag(1, "tag1", 0),
				testutil.MustNewTag(2, "tag2", 1),
				testutil.MustNewTag(3, "tag3", 2),
			}),
			args:      args{idList: []model.TagID{4}},
			wantM:     map[model.TagID]*model.Tag{},
			wantNameM: map[string]*model.Tag{},
		},
		{
			name: "return subset with two tags",
			tagSet: model.NewTagSet([]*model.Tag{
				testutil.MustNewTag(1, "tag1", 0),
				testutil.MustNewTag(2, "tag2", 1),
				testutil.MustNewTag(3, "tag3", 2),
			}),
			args: args{idList: []model.TagID{1, 2}},
			wantM: map[model.TagID]*model.Tag{
				1: testutil.MustNewTag(1, "tag1", 0),
				2: testutil.MustNewTag(2, "tag2", 1),
			},
			wantNameM: map[string]*model.Tag{
				"tag1": testutil.MustNewTag(1, "tag1", 0),
				"tag2": testutil.MustNewTag(2, "tag2", 1),
			},
		},
		{
			name: "return subset with one tag",
			tagSet: model.NewTagSet([]*model.Tag{
				testutil.MustNewTag(1, "tag1", 0),
				testutil.MustNewTag(2, "tag2", 1),
				testutil.MustNewTag(3, "tag3", 2),
			}),
			args: args{idList: []model.TagID{1}},
			wantM: map[model.TagID]*model.Tag{
				1: testutil.MustNewTag(1, "tag1", 0),
			},
			wantNameM: map[string]*model.Tag{
				"tag1": testutil.MustNewTag(1, "tag1", 0),
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := tt.tagSet
			subset := t.SubSetByID(tt.args.idList)
			gotM, gotNameM := subset.ToMap()
			testutil.Diff(t1, tt.wantM, gotM)
			testutil.Diff(t1, tt.wantNameM, gotNameM)
		})
	}
}

func TestTagSet_SplitByID(t1 *testing.T) {
	type args struct {
		idList []model.TagID
	}
	tests := []struct {
		name       string
		tagSet     *model.TagSet
		args       args
		wantTrueM  map[model.TagID]*model.Tag
		wantFalseM map[model.TagID]*model.Tag
	}{
		{
			name:       "return empty subset if TagSet is empty",
			tagSet:     model.NewTagSet([]*model.Tag{}),
			args:       args{idList: []model.TagID{1}},
			wantTrueM:  map[model.TagID]*model.Tag{},
			wantFalseM: map[model.TagID]*model.Tag{},
		},
		{
			name: "all tags in false tag set if non exists ID is given",
			tagSet: model.NewTagSet([]*model.Tag{
				testutil.MustNewTag(1, "tag1", 0),
				testutil.MustNewTag(2, "tag2", 1),
				testutil.MustNewTag(3, "tag3", 2),
			}),
			args:      args{idList: []model.TagID{4}},
			wantTrueM: map[model.TagID]*model.Tag{},
			wantFalseM: map[model.TagID]*model.Tag{
				1: testutil.MustNewTag(1, "tag1", 0),
				2: testutil.MustNewTag(2, "tag2", 1),
				3: testutil.MustNewTag(3, "tag3", 2),
			},
		},
		{
			name: "two match",
			tagSet: model.NewTagSet([]*model.Tag{
				testutil.MustNewTag(1, "tag1", 0),
				testutil.MustNewTag(2, "tag2", 1),
				testutil.MustNewTag(3, "tag3", 2),
			}),
			args: args{idList: []model.TagID{1, 2}},
			wantTrueM: map[model.TagID]*model.Tag{
				1: testutil.MustNewTag(1, "tag1", 0),
				2: testutil.MustNewTag(2, "tag2", 1),
			},
			wantFalseM: map[model.TagID]*model.Tag{
				3: testutil.MustNewTag(3, "tag3", 2),
			},
		},
		{
			name: "one true",
			tagSet: model.NewTagSet([]*model.Tag{
				testutil.MustNewTag(1, "tag1", 0),
				testutil.MustNewTag(2, "tag2", 1),
				testutil.MustNewTag(3, "tag3", 2),
			}),
			args: args{idList: []model.TagID{1}},
			wantTrueM: map[model.TagID]*model.Tag{
				1: testutil.MustNewTag(1, "tag1", 0),
			},
			wantFalseM: map[model.TagID]*model.Tag{
				2: testutil.MustNewTag(2, "tag2", 1),
				3: testutil.MustNewTag(3, "tag3", 2),
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := tt.tagSet
			trueSet, falseSet := t.SplitByID(tt.args.idList)
			trueM, _ := trueSet.ToMap()
			falseM, _ := falseSet.ToMap()
			testutil.Diff(t1, tt.wantTrueM, trueM)
			testutil.Diff(t1, tt.wantFalseM, falseM)
		})
	}
}

func TestTagSet_SubSetByNames(t1 *testing.T) {
	type args struct {
		names []string
	}
	tests := []struct {
		name      string
		tagSet    *model.TagSet
		args      args
		wantM     map[model.TagID]*model.Tag
		wantNameM map[string]*model.Tag
	}{
		{
			name:      "return empty subset if TagSet is empty",
			tagSet:    model.NewTagSet([]*model.Tag{}),
			args:      args{names: []string{"tag1"}},
			wantM:     map[model.TagID]*model.Tag{},
			wantNameM: map[string]*model.Tag{},
		},
		{
			name: "return empty subset if any tag does not match",
			tagSet: model.NewTagSet([]*model.Tag{
				testutil.MustNewTag(1, "tag1", 0),
				testutil.MustNewTag(2, "tag2", 1),
				testutil.MustNewTag(3, "tag3", 2),
			}),
			args:      args{names: []string{"tag4"}},
			wantM:     map[model.TagID]*model.Tag{},
			wantNameM: map[string]*model.Tag{},
		},
		{
			name: "return subset with two tags",
			tagSet: model.NewTagSet([]*model.Tag{
				testutil.MustNewTag(1, "tag1", 0),
				testutil.MustNewTag(2, "tag2", 1),
				testutil.MustNewTag(3, "tag3", 2),
			}),
			args: args{names: []string{"tag1", "tag2"}},
			wantM: map[model.TagID]*model.Tag{
				1: testutil.MustNewTag(1, "tag1", 0),
				2: testutil.MustNewTag(2, "tag2", 1),
			},
			wantNameM: map[string]*model.Tag{
				"tag1": testutil.MustNewTag(1, "tag1", 0),
				"tag2": testutil.MustNewTag(2, "tag2", 1),
			},
		},
		{
			name: "return subset with one tag",
			tagSet: model.NewTagSet([]*model.Tag{
				testutil.MustNewTag(1, "tag1", 0),
				testutil.MustNewTag(2, "tag2", 1),
				testutil.MustNewTag(3, "tag3", 2),
			}),
			args: args{names: []string{"tag1"}},
			wantM: map[model.TagID]*model.Tag{
				1: testutil.MustNewTag(1, "tag1", 0),
			},
			wantNameM: map[string]*model.Tag{
				"tag1": testutil.MustNewTag(1, "tag1", 0),
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := tt.tagSet
			subset := t.SubSetByNames(tt.args.names)
			gotM, gotNameM := subset.ToMap()
			testutil.Diff(t1, tt.wantM, gotM)
			testutil.Diff(t1, tt.wantNameM, gotNameM)
		})
	}
}

func TestTagSet_SplitByNames(t1 *testing.T) {
	type args struct {
		tagNames []string
	}
	tests := []struct {
		name       string
		tagSet     *model.TagSet
		args       args
		wantTrueM  map[model.TagID]*model.Tag
		wantFalseM map[model.TagID]*model.Tag
	}{
		{
			name:       "return empty subset if TagSet is empty",
			tagSet:     model.NewTagSet([]*model.Tag{}),
			args:       args{tagNames: []string{"tag1"}},
			wantTrueM:  map[model.TagID]*model.Tag{},
			wantFalseM: map[model.TagID]*model.Tag{},
		},
		{
			name: "all tags in false tag set if non exists ID is given",
			tagSet: model.NewTagSet([]*model.Tag{
				testutil.MustNewTag(1, "tag1", 0),
				testutil.MustNewTag(2, "tag2", 1),
				testutil.MustNewTag(3, "tag3", 2),
			}),
			args:      args{tagNames: []string{"tag4"}},
			wantTrueM: map[model.TagID]*model.Tag{},
			wantFalseM: map[model.TagID]*model.Tag{
				1: testutil.MustNewTag(1, "tag1", 0),
				2: testutil.MustNewTag(2, "tag2", 1),
				3: testutil.MustNewTag(3, "tag3", 2),
			},
		},
		{
			name: "two match",
			tagSet: model.NewTagSet([]*model.Tag{
				testutil.MustNewTag(1, "tag1", 0),
				testutil.MustNewTag(2, "tag2", 1),
				testutil.MustNewTag(3, "tag3", 2),
			}),
			args: args{tagNames: []string{"tag1", "tag2"}},
			wantTrueM: map[model.TagID]*model.Tag{
				1: testutil.MustNewTag(1, "tag1", 0),
				2: testutil.MustNewTag(2, "tag2", 1),
			},
			wantFalseM: map[model.TagID]*model.Tag{
				3: testutil.MustNewTag(3, "tag3", 2),
			},
		},
		{
			name: "one true",
			tagSet: model.NewTagSet([]*model.Tag{
				testutil.MustNewTag(1, "tag1", 0),
				testutil.MustNewTag(2, "tag2", 1),
				testutil.MustNewTag(3, "tag3", 2),
			}),
			args: args{tagNames: []string{"tag1"}},
			wantTrueM: map[model.TagID]*model.Tag{
				1: testutil.MustNewTag(1, "tag1", 0),
			},
			wantFalseM: map[model.TagID]*model.Tag{
				2: testutil.MustNewTag(2, "tag2", 1),
				3: testutil.MustNewTag(3, "tag3", 2),
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := tt.tagSet
			trueSet, falseSet := t.SplitByNames(tt.args.tagNames)
			trueM, _ := trueSet.ToMap()
			falseM, _ := falseSet.ToMap()
			testutil.Diff(t1, tt.wantTrueM, trueM)
			testutil.Diff(t1, tt.wantFalseM, falseM)
		})
	}
}

func TestTagSet_ToMap(t1 *testing.T) {
	tests := []struct {
		name      string
		tagSet    *model.TagSet
		wantM     map[model.TagID]*model.Tag
		wantNameM map[string]*model.Tag
	}{
		{
			name:      "return empty maps if TagSet is empty",
			tagSet:    model.NewTagSet([]*model.Tag{}),
			wantM:     map[model.TagID]*model.Tag{},
			wantNameM: map[string]*model.Tag{},
		},
		{
			name: "return subset with two tags",
			tagSet: model.NewTagSet([]*model.Tag{
				testutil.MustNewTag(1, "tag1", 0),
				testutil.MustNewTag(2, "tag2", 1),
				testutil.MustNewTag(3, "tag3", 2),
			}),
			wantM: map[model.TagID]*model.Tag{
				1: testutil.MustNewTag(1, "tag1", 0),
				2: testutil.MustNewTag(2, "tag2", 1),
				3: testutil.MustNewTag(3, "tag3", 2),
			},
			wantNameM: map[string]*model.Tag{
				"tag1": testutil.MustNewTag(1, "tag1", 0),
				"tag2": testutil.MustNewTag(2, "tag2", 1),
				"tag3": testutil.MustNewTag(3, "tag3", 2),
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			gotM, gotNameM := tt.tagSet.ToMap()
			testutil.Diff(t1, tt.wantM, gotM)
			testutil.Diff(t1, tt.wantNameM, gotNameM)
		})
	}
}

func TestTag_GetID(t1 *testing.T) {
	type fields struct {
		ID   model.TagID
		Name string
	}
	tests := []struct {
		name   string
		fields fields
		want   uint64
	}{
		{
			name:   "return ID",
			fields: fields{ID: 1, Name: "tag1"},
			want:   1,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			tag := &model.UnindexedTag{
				ID:   tt.fields.ID,
				Name: tt.fields.Name,
			}
			testutil.Diff(t1, tt.want, tag.GetID())
		})
	}
}

func TestTag_SetID(t1 *testing.T) {
	type args struct {
		id uint64
	}
	tests := []struct {
		name         string
		unindexedTag *model.UnindexedTag
		args         args
	}{
		{unindexedTag: &model.UnindexedTag{ID: 1, Name: "tag1"}, args: args{id: 2}},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			tt.unindexedTag.SetID(tt.args.id)
			testutil.Diff(t1, tt.args.id, uint64(tt.unindexedTag.ID))
		})
	}
}
