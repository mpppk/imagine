package model

import (
	"testing"

	"github.com/mpppk/imagine/testutil"
)

func TestNewTagSet(t *testing.T) {
	type args struct {
		tags []*Tag
	}
	tests := []struct {
		name      string
		args      args
		wantM     map[TagID]*Tag
		wantNameM map[string]*Tag
	}{
		{
			name:      "return empty TagSet from empty tags",
			args:      args{tags: []*Tag{}},
			wantM:     map[TagID]*Tag{},
			wantNameM: map[string]*Tag{},
		},
		{
			name:      "return TagSet from one tag",
			args:      args{tags: []*Tag{{ID: 1, Name: "tag1"}}},
			wantM:     map[TagID]*Tag{1: {ID: 1, Name: "tag1"}},
			wantNameM: map[string]*Tag{"tag1": {ID: 1, Name: "tag1"}},
		},
		{
			name:      "return TagSet from two tags",
			args:      args{tags: []*Tag{{ID: 1, Name: "tag1"}, {ID: 2, Name: "tag2"}}},
			wantM:     map[TagID]*Tag{1: {ID: 1, Name: "tag1"}, 2: {ID: 2, Name: "tag2"}},
			wantNameM: map[string]*Tag{"tag1": {ID: 1, Name: "tag1"}, "tag2": {ID: 2, Name: "tag2"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tagSet := NewTagSet(tt.args.tags)
			m, nameM := tagSet.ToMap()
			testutil.Diff(t, tt.wantM, m)
			testutil.Diff(t, tt.wantNameM, nameM)
		})
	}
}

func TestTagSet_Get(t1 *testing.T) {
	type args struct {
		id TagID
	}
	tests := []struct {
		name   string
		tagSet *TagSet
		args   args
		want   *Tag
		want1  bool
	}{
		{
			name:   "get from empty TagSet",
			tagSet: NewTagSet([]*Tag{}),
			args:   args{id: 1},
			want:   nil,
			want1:  false,
		},
		{
			name:   "get exist tag from one tag",
			tagSet: NewTagSet([]*Tag{{ID: 1, Name: "tag1"}}),
			args:   args{id: 1},
			want:   &Tag{ID: 1, Name: "tag1"},
			want1:  true,
		},
		{
			name:   "get non exist tag from one tag",
			tagSet: NewTagSet([]*Tag{{ID: 1, Name: "tag1"}}),
			args:   args{id: 2},
			want:   nil,
			want1:  false,
		},
		{
			name:   "get exist tag from two tags",
			tagSet: NewTagSet([]*Tag{{ID: 1, Name: "tag1"}, {ID: 2, Name: "tag2"}}),
			args:   args{id: 2},
			want:   &Tag{ID: 2, Name: "tag2"},
			want1:  true,
		},
		{
			name:   "get exist tag from two tags",
			tagSet: NewTagSet([]*Tag{{ID: 1, Name: "tag1"}, {ID: 2, Name: "tag2"}}),
			args:   args{id: 3},
			want:   nil,
			want1:  false,
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
		tagSet *TagSet
		args   args
		want   *Tag
		want1  bool
	}{
		{
			name:   "get from empty TagSet",
			tagSet: NewTagSet([]*Tag{}),
			args:   args{name: "tag1"},
			want:   nil,
			want1:  false,
		},
		{
			name:   "get exist tag from one tag",
			tagSet: NewTagSet([]*Tag{{ID: 1, Name: "tag1"}}),
			args:   args{name: "tag1"},
			want:   &Tag{ID: 1, Name: "tag1"},
			want1:  true,
		},
		{
			name:   "get non exist tag from one tag",
			tagSet: NewTagSet([]*Tag{{ID: 1, Name: "tag1"}}),
			args:   args{name: "tag2"},
			want:   nil,
			want1:  false,
		},
		{
			name:   "get exist tag from two tags",
			tagSet: NewTagSet([]*Tag{{ID: 1, Name: "tag1"}, {ID: 2, Name: "tag2"}}),
			args:   args{name: "tag2"},
			want:   &Tag{ID: 2, Name: "tag2"},
			want1:  true,
		},
		{
			name:   "get exist tag from two tags",
			tagSet: NewTagSet([]*Tag{{ID: 1, Name: "tag1"}, {ID: 2, Name: "tag2"}}),
			args:   args{name: "tag3"},
			want:   nil,
			want1:  false,
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
		tag *Tag
	}
	tests := []struct {
		name      string
		tagSet    *TagSet
		args      args
		want      bool
		wantM     map[TagID]*Tag
		wantNameM map[string]*Tag
	}{
		{
			name:      "add tag to TagSet",
			tagSet:    NewTagSet([]*Tag{{ID: 1, Name: "tag1"}, {ID: 2, Name: "tag2"}}),
			args:      args{tag: &Tag{ID: 3, Name: "tag3"}},
			want:      true,
			wantM:     map[TagID]*Tag{1: {ID: 1, Name: "tag1"}, 2: {ID: 2, Name: "tag2"}, 3: {ID: 3, Name: "tag3"}},
			wantNameM: map[string]*Tag{"tag1": {ID: 1, Name: "tag1"}, "tag2": {ID: 2, Name: "tag2"}, "tag3": {ID: 3, Name: "tag3"}},
		},
		{
			name:      "update tag",
			tagSet:    NewTagSet([]*Tag{{ID: 1, Name: "tag1"}, {ID: 2, Name: "tag2"}}),
			args:      args{tag: &Tag{ID: 2, Name: "updated-tag2"}},
			want:      true,
			wantM:     map[TagID]*Tag{1: {ID: 1, Name: "tag1"}, 2: {ID: 2, Name: "updated-tag2"}},
			wantNameM: map[string]*Tag{"tag1": {ID: 1, Name: "tag1"}, "tag2": {ID: 2, Name: "updated-tag2"}},
		},
		{
			name:      "do nothing for same tag",
			tagSet:    NewTagSet([]*Tag{{ID: 1, Name: "tag1"}, {ID: 2, Name: "tag2"}}),
			args:      args{tag: &Tag{ID: 2, Name: "tag2"}},
			want:      true,
			wantM:     map[TagID]*Tag{1: {ID: 1, Name: "tag1"}, 2: {ID: 2, Name: "tag2"}},
			wantNameM: map[string]*Tag{"tag1": {ID: 1, Name: "tag1"}, "tag2": {ID: 2, Name: "tag2"}},
		},
		{
			name:      "do nothing to tag which have same name but different ID, then return false",
			tagSet:    NewTagSet([]*Tag{{ID: 1, Name: "tag1"}, {ID: 2, Name: "tag2"}}),
			args:      args{tag: &Tag{ID: 3, Name: "tag2"}},
			want:      false,
			wantM:     map[TagID]*Tag{1: {ID: 1, Name: "tag1"}, 2: {ID: 2, Name: "tag2"}},
			wantNameM: map[string]*Tag{"tag1": {ID: 1, Name: "tag1"}, "tag2": {ID: 2, Name: "tag2"}},
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
		f func(tag *Tag) bool
	}
	tests := []struct {
		name      string
		tagSet    *TagSet
		args      args
		wantM     map[TagID]*Tag
		wantNameM map[string]*Tag
	}{
		{
			name:      "return empty subset if TagSet is empty",
			tagSet:    NewTagSet([]*Tag{}),
			args:      args{f: func(tag *Tag) bool { return true }},
			wantM:     map[TagID]*Tag{},
			wantNameM: map[string]*Tag{},
		},
		{
			name:      "return empty subset if any tag does not match",
			tagSet:    NewTagSet([]*Tag{{ID: 1, Name: "tag1"}, {ID: 2, Name: "tag2"}, {ID: 3, Name: "tag3"}}),
			args:      args{f: func(tag *Tag) bool { return false }},
			wantM:     map[TagID]*Tag{},
			wantNameM: map[string]*Tag{},
		},
		{
			name:   "return subset with two tags",
			tagSet: NewTagSet([]*Tag{{ID: 1, Name: "tag1"}, {ID: 2, Name: "tag2"}, {ID: 3, Name: "tag3"}}),
			args: args{f: func(tag *Tag) bool {
				return tag.Name == "tag1" || tag.Name == "tag2"
			}},
			wantM:     map[TagID]*Tag{1: {ID: 1, Name: "tag1"}, 2: {ID: 2, Name: "tag2"}},
			wantNameM: map[string]*Tag{"tag1": {ID: 1, Name: "tag1"}, "tag2": {ID: 2, Name: "tag2"}},
		},
		{
			name:   "return subset with one tag",
			tagSet: NewTagSet([]*Tag{{ID: 1, Name: "tag1"}, {ID: 2, Name: "tag2"}, {ID: 3, Name: "tag3"}}),
			args: args{f: func(tag *Tag) bool {
				return tag.Name == "tag1"
			}},
			wantM:     map[TagID]*Tag{1: {ID: 1, Name: "tag1"}},
			wantNameM: map[string]*Tag{"tag1": {ID: 1, Name: "tag1"}},
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
		f func(tag *Tag) bool
	}
	tests := []struct {
		name       string
		tagSet     *TagSet
		args       args
		wantTrueM  map[TagID]*Tag
		wantFalseM map[TagID]*Tag
	}{
		{
			name:       "return empty subset if TagSet is empty",
			tagSet:     NewTagSet([]*Tag{}),
			args:       args{f: func(tag *Tag) bool { return true }},
			wantTrueM:  map[TagID]*Tag{},
			wantFalseM: map[TagID]*Tag{},
		},
		{
			name:       "all false",
			tagSet:     NewTagSet([]*Tag{{ID: 1, Name: "tag1"}, {ID: 2, Name: "tag2"}, {ID: 3, Name: "tag3"}}),
			args:       args{f: func(tag *Tag) bool { return false }},
			wantTrueM:  map[TagID]*Tag{},
			wantFalseM: map[TagID]*Tag{1: {ID: 1, Name: "tag1"}, 2: {ID: 2, Name: "tag2"}, 3: {ID: 3, Name: "tag3"}},
		},
		{
			name:   "two true",
			tagSet: NewTagSet([]*Tag{{ID: 1, Name: "tag1"}, {ID: 2, Name: "tag2"}, {ID: 3, Name: "tag3"}}),
			args: args{f: func(tag *Tag) bool {
				return tag.Name == "tag1" || tag.Name == "tag2"
			}},
			wantTrueM:  map[TagID]*Tag{1: {ID: 1, Name: "tag1"}, 2: {ID: 2, Name: "tag2"}},
			wantFalseM: map[TagID]*Tag{3: {ID: 3, Name: "tag3"}},
		},
		{
			name:   "one true",
			tagSet: NewTagSet([]*Tag{{ID: 1, Name: "tag1"}, {ID: 2, Name: "tag2"}, {ID: 3, Name: "tag3"}}),
			args: args{f: func(tag *Tag) bool {
				return tag.Name == "tag1"
			}},
			wantTrueM:  map[TagID]*Tag{1: {ID: 1, Name: "tag1"}},
			wantFalseM: map[TagID]*Tag{2: {ID: 2, Name: "tag2"}, 3: {ID: 3, Name: "tag3"}},
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

func TestTagSet_SubSetByNames(t1 *testing.T) {
	type args struct {
		names []string
	}
	tests := []struct {
		name      string
		tagSet    *TagSet
		args      args
		wantM     map[TagID]*Tag
		wantNameM map[string]*Tag
	}{
		{
			name:      "return empty subset if TagSet is empty",
			tagSet:    NewTagSet([]*Tag{}),
			args:      args{names: []string{"tag1"}},
			wantM:     map[TagID]*Tag{},
			wantNameM: map[string]*Tag{},
		},
		{
			name:      "return empty subset if any tag does not match",
			tagSet:    NewTagSet([]*Tag{{ID: 1, Name: "tag1"}, {ID: 2, Name: "tag2"}, {ID: 3, Name: "tag3"}}),
			args:      args{names: []string{"tag4"}},
			wantM:     map[TagID]*Tag{},
			wantNameM: map[string]*Tag{},
		},
		{
			name:      "return subset with two tags",
			tagSet:    NewTagSet([]*Tag{{ID: 1, Name: "tag1"}, {ID: 2, Name: "tag2"}, {ID: 3, Name: "tag3"}}),
			args:      args{names: []string{"tag1", "tag2"}},
			wantM:     map[TagID]*Tag{1: {ID: 1, Name: "tag1"}, 2: {ID: 2, Name: "tag2"}},
			wantNameM: map[string]*Tag{"tag1": {ID: 1, Name: "tag1"}, "tag2": {ID: 2, Name: "tag2"}},
		},
		{
			name:      "return subset with one tag",
			tagSet:    NewTagSet([]*Tag{{ID: 1, Name: "tag1"}, {ID: 2, Name: "tag2"}, {ID: 3, Name: "tag3"}}),
			args:      args{names: []string{"tag1"}},
			wantM:     map[TagID]*Tag{1: {ID: 1, Name: "tag1"}},
			wantNameM: map[string]*Tag{"tag1": {ID: 1, Name: "tag1"}},
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

func TestTagSet_ToMap(t1 *testing.T) {
	tests := []struct {
		name      string
		tagSet    *TagSet
		wantM     map[TagID]*Tag
		wantNameM map[string]*Tag
	}{
		{
			name:      "return empty maps if TagSet is empty",
			tagSet:    NewTagSet([]*Tag{}),
			wantM:     map[TagID]*Tag{},
			wantNameM: map[string]*Tag{},
		},
		{
			name:      "return subset with two tags",
			tagSet:    NewTagSet([]*Tag{{ID: 1, Name: "tag1"}, {ID: 2, Name: "tag2"}, {ID: 3, Name: "tag3"}}),
			wantM:     map[TagID]*Tag{1: {ID: 1, Name: "tag1"}, 2: {ID: 2, Name: "tag2"}, 3: {ID: 3, Name: "tag3"}},
			wantNameM: map[string]*Tag{"tag1": {ID: 1, Name: "tag1"}, "tag2": {ID: 2, Name: "tag2"}, "tag3": {ID: 3, Name: "tag3"}},
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
		ID   TagID
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
			tag := &Tag{
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
		name string
		tag  *Tag
		args args
	}{
		{tag: &Tag{ID: 1, Name: "tag1"}, args: args{id: 2}},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			tt.tag.SetID(tt.args.id)
			testutil.Diff(t1, tt.args.id, uint64(tt.tag.ID))
		})
	}
}