package action

import (
	"testing"

	"github.com/mpppk/imagine/domain/model"

	"github.com/mpppk/imagine/usecase/mock_usecase"

	"github.com/golang/mock/gomock"

	fsa "github.com/mpppk/lorca-fsa/lorca-fsa"
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
				Tags:      []*model.Tag{},
			}},
			wantActions: []*fsa.Action{tagActionCreator.save("default-workspace", []*model.Tag{})},
			wantErr:     false,
		},
		{
			name: "Save one tag",
			args: args{payload: &tagUpdatePayload{
				WsPayload: WsPayload{"default-workspace"},
				Tags:      []*model.Tag{{Name: "tag1"}},
			}},
			wantActions: []*fsa.Action{tagActionCreator.save("default-workspace", []*model.Tag{{Name: "tag1"}})},
			wantErr:     false,
		},
		{
			name: "Save two tags",
			args: args{payload: &tagUpdatePayload{
				WsPayload: WsPayload{"default-workspace"},
				Tags:      []*model.Tag{{Name: "tag1"}, {Name: "tag2"}},
			}},
			wantActions: []*fsa.Action{tagActionCreator.save("default-workspace", []*model.Tag{{Name: "tag1"}, {Name: "tag2"}})},
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			dispatcher := newMockDispatcher(t1, tt.wantActions)
			defer dispatcher.Finish()
			ctrl := gomock.NewController(t1)
			tagUseCase := mock_usecase.NewMockTag(ctrl)
			tagUseCase.EXPECT().PutTags(gomock.Eq(tt.args.payload.WorkSpaceName), gomock.Eq(tt.args.payload.Tags))

			t := &tagSaveHandler{
				tagUseCase: tagUseCase,
				action:     tagActionCreator,
			}

			action := &fsa.Action{Payload: tt.args.payload}
			if err := t.Do(action, dispatcher.Dispatch); (err != nil) != tt.wantErr {
				t1.Errorf("Do() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
