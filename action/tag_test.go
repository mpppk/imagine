package action

import (
	"testing"

	"github.com/mpppk/imagine/domain/service/tagsvc"

	"github.com/mitchellh/mapstructure"

	"github.com/mpppk/imagine/testutil"

	"github.com/mpppk/imagine/domain/model"

	"github.com/mpppk/imagine/usecase/mock_usecase"

	"github.com/golang/mock/gomock"

	fsa "github.com/mpppk/lorca-fsa"
)

func Test_tagSaveHandler_Do(t1 *testing.T) {
	tagActionCreator := &tagActionCreator{}
	type args struct {
		payload *tagUpdatePayload
	}
	tests := []struct {
		name        string
		args        args
		wantActions []*fsa.Action
		wantErr     bool
	}{
		{
			name: "Save empty tags",
			args: args{payload: &tagUpdatePayload{
				WsPayload: WsPayload{"default-workspace"},
				Tags:      []*model.UnindexedTag{},
			}},
			wantActions: []*fsa.Action{tagActionCreator.save("default-workspace", nil)},
			wantErr:     false,
		},
		{
			name: "Save one tag",
			args: args{payload: &tagUpdatePayload{
				WsPayload: WsPayload{"default-workspace"},
				Tags:      []*model.UnindexedTag{testutil.MustNewUnindexedTag(1, "tag1")},
			}},
			wantActions: []*fsa.Action{tagActionCreator.save("default-workspace", []*model.UnindexedTag{
				testutil.MustNewUnindexedTag(1, "tag1"),
			})},
			wantErr: false,
		},
		{
			name: "Save two tags",
			args: args{payload: &tagUpdatePayload{
				WsPayload: WsPayload{"default-workspace"},
				Tags: []*model.UnindexedTag{
					testutil.MustNewUnindexedTag(1, "tag1"),
					testutil.MustNewUnindexedTag(2, "tag2"),
				},
			}},
			wantActions: []*fsa.Action{tagActionCreator.save("default-workspace", []*model.UnindexedTag{
				testutil.MustNewUnindexedTag(1, "tag1"),
				testutil.MustNewUnindexedTag(2, "tag2"),
			})},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			dispatcher := testutil.NewMockDispatcher(t1)
			ctrl := gomock.NewController(t1)
			tagUseCase := mock_usecase.NewMockTag(ctrl)
			var setTagsRet []*model.Tag = nil
			if len(tt.wantActions) > 0 {
				var payload tagSavePayload
				if err := mapstructure.Decode(tt.wantActions[0].Payload, &payload); err != nil {
					t1.Fatalf("failed to parse action payload: %#v", tt.wantActions[0])
				}
				setTagsRet = tagsvc.ToTagsWithOrderIndex(payload.Tags)
			}
			tagUseCase.EXPECT().SetTags(gomock.Eq(tt.args.payload.WorkSpaceName), gomock.Eq(tt.args.payload.Tags)).Return(setTagsRet, nil)

			t := &tagSaveHandler{
				tagUseCase: tagUseCase,
				action:     tagActionCreator,
			}

			action := &fsa.Action{Payload: tt.args.payload}
			if err := t.Do(action, dispatcher.Dispatch); (err != nil) != tt.wantErr {
				t1.Errorf("Do() error = %v, wantErr %v", err, tt.wantErr)
			}
			dispatcher.Test(tt.wantActions)
		})
	}
}
