package repoimpl_test

import (
	"testing"

	"github.com/mpppk/imagine/testutil"

	"github.com/mpppk/imagine/domain/model"
	"github.com/mpppk/imagine/usecase/usecasetest"
)

func TestBBoltTag_SetTags(t *testing.T) {
	type args struct {
		ws       model.WSName
		tagNames []string
	}
	tests := []struct {
		name          string
		args1         args
		args2         args
		existTagNames []string
		want1         []*model.TagWithIndex
		want2         []*model.TagWithIndex
		wantTags1     []*model.TagWithIndex
		wantTags2     []*model.TagWithIndex
		wantErr       bool
	}{
		{
			name:          "add two tags to empty db",
			existTagNames: []string{"tag1"},
			args1: args{
				ws:       "default-workspace",
				tagNames: []string{"tag1", "tag2"},
			},
			want1: []*model.TagWithIndex{
				testutil.MustNewTagWithIndex(1, "tag1", 0),
				testutil.MustNewTagWithIndex(2, "tag2", 1),
			},
			wantTags1: []*model.TagWithIndex{
				testutil.MustNewTagWithIndex(1, "tag1", 0),
				testutil.MustNewTagWithIndex(2, "tag2", 1),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		usecasetest.RunParallelWithUseCases(t, tt.name, tt.args1.ws, func(t *testing.T, ut *usecasetest.UseCases) {
			f := func(args args, want, wantTags []*model.TagWithIndex) {
				tags, err := ut.Usecases.Tag.SetTagByNames(args.ws, args.tagNames)
				if testutil.FatalIfErrIsUnexpected(t, tt.wantErr, err) {
					return
				}
				testutil.Diff(t, want, tags)

				newTags := ut.Tag.List(args.ws)
				testutil.Diff(t, wantTags, newTags)
			}
			f(tt.args1, tt.want1, tt.wantTags1)
			if tt.want2 != nil {
				f(tt.args2, tt.want2, tt.wantTags2)
			}
		})
	}
}

func TestBBoltTag_Save(t *testing.T) {
	type args struct {
		ws  model.WSName
		tag *model.TagWithIndex
	}
	tests := []struct {
		name          string
		args          args
		existTagNames []string
		want          *model.TagWithIndex
		wantTags      []*model.Tag
		wantErr       bool
	}{
		{
			name:          "update tag",
			existTagNames: []string{"tag1"},
			args: args{
				ws:  "default-workspace",
				tag: testutil.MustNewTagWithIndex(1, "updated-tag1", 0),
			},
			want:     testutil.MustNewTagWithIndex(1, "updated-tag1", 0),
			wantTags: []*model.Tag{{ID: 1, Name: "updated-tag1"}},
			wantErr:  false,
		},
		{
			name:          "add new tag",
			existTagNames: []string{"tag1"},
			args: args{
				ws:  "default-workspace",
				tag: testutil.MustNewTagWithIndex(0, "updated-tag2", 1),
			},
			want:     testutil.MustNewTagWithIndex(2, "updated-tag2", 1),
			wantTags: []*model.Tag{{ID: 1, Name: "updated-tag1"}, {ID: 2, Name: "updated-tag2"}},
			wantErr:  false,
		},
		{
			name:          "fail if tag ID which does not exist is provided",
			existTagNames: []string{"tag1"},
			args: args{
				ws:  "default-workspace",
				tag: testutil.MustNewTagWithIndex(2, "updated-tag2", 1),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		usecasetest.RunParallelWithUseCases(t, tt.name, tt.args.ws, func(t *testing.T, ut *usecasetest.UseCases) {
			ut.Tag.SetTagByNames(tt.args.ws, tt.existTagNames)

			id, err := ut.Usecases.Client.Tag.Save(tt.args.ws, tt.args.tag)
			if testutil.FatalIfErrIsUnexpected(t, tt.wantErr, err) {
				return
			}
			testutil.Diff(t, tt.want, id)

			got, _, err := ut.Usecases.Client.Tag.Get(tt.args.ws, tt.args.tag.ID)
			if err != nil {
				t.Errorf("failed to get tag: %v: %v", tt.args.tag.ID, err)
			}

			testutil.Diff(t, tt.args.tag, got)
		})
	}
}

func TestBBoltTag_Add(t *testing.T) {
	type args struct {
		tag *model.UnregisteredTag
	}
	tests := []struct {
		name     string
		ws       model.WSName
		argsList []args
		wants    []*model.TagWithIndex
		wantTags []*model.TagWithIndex
		wantErr  bool
	}{
		{
			name: "add two tags",
			ws:   "default-workspace",
			argsList: []args{
				{tag: testutil.MustNewUnregisteredTag("tag1")},
				{tag: testutil.MustNewUnregisteredTag("tag2")},
			},
			wants: []*model.TagWithIndex{
				testutil.MustNewTagWithIndex(1, "tag1", 0),
				testutil.MustNewTagWithIndex(2, "tag2", 1),
			},
			wantTags: []*model.TagWithIndex{
				testutil.MustNewTagWithIndex(1, "tag1", 0),
				testutil.MustNewTagWithIndex(2, "tag2", 1),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		usecasetest.RunParallelWithUseCases(t, tt.name, tt.ws, func(t *testing.T, ut *usecasetest.UseCases) {
			for i, args := range tt.argsList {
				tag, err := ut.Usecases.Client.Tag.Add(tt.ws, args.tag)
				if (err != nil) != tt.wantErr {
					t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				testutil.Diff(t, tt.wants[i], tag)
			}
			gotTags := ut.Tag.List(tt.ws)
			testutil.Diff(t, tt.wantTags, gotTags)
		})
	}
}

func TestBBoltTag_AddByNames(t *testing.T) {
	type args struct {
		ws       model.WSName
		tagNames []string
	}
	tests := []struct {
		name          string
		args          args
		existTagNames []string
		want          []*model.TagWithIndex
		wantTags      []*model.TagWithIndex
		wantErr       bool
	}{
		{
			name:          "add tag",
			existTagNames: []string{"tag1"},
			args: args{
				ws:       "default-workspace",
				tagNames: []string{"tag2"},
			},
			want: []*model.TagWithIndex{
				testutil.MustNewTagWithIndex(2, "tag2", 1),
			},
			wantTags: []*model.TagWithIndex{
				testutil.MustNewTagWithIndex(1, "tag1", 0),
				testutil.MustNewTagWithIndex(2, "tag2", 1),
			},
			wantErr: false,
		},
		{
			name:          "add tag to empty DB",
			existTagNames: []string{"tag1"},
			args: args{
				ws:       "default-workspace",
				tagNames: []string{"tag1"},
			},
			want: []*model.TagWithIndex{
				testutil.MustNewTagWithIndex(1, "tag1", 0),
			},
			wantTags: []*model.TagWithIndex{
				testutil.MustNewTagWithIndex(1, "tag1", 0),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		usecasetest.RunParallelWithUseCases(t, tt.name, tt.args.ws, func(t *testing.T, ut *usecasetest.UseCases) {
			ut.Tag.SetTagByNames(tt.args.ws, tt.existTagNames)

			idList, err := ut.Usecases.Client.Tag.AddByNames(tt.args.ws, tt.args.tagNames)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddByNames() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			testutil.Diff(t, tt.want, idList)

			gotTags := ut.Tag.List(tt.args.ws)
			testutil.Diff(t, tt.wantTags, gotTags)
		})
	}
}
