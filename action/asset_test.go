package action

import (
	"testing"

	"github.com/mpppk/imagine/usecase/mock_usecase"

	"github.com/golang/mock/gomock"
	"github.com/mpppk/imagine/testutil"

	fsa "github.com/mpppk/lorca-fsa"
	lorcafsa "github.com/mpppk/lorca-fsa"
)

func Test_assetScanHandler_Do(t *testing.T) {
	type args struct {
		action   *fsa.Action
		dispatch lorcafsa.Dispatch
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		dispatcher := testutil.NewMockDispatcher(t1, tt.wantActions)
		defer dispatcher.Finish()
		ctrl := gomock.NewController(t1)
		assetUseCase := mock_usecase.MockTag{}
		t.Run(tt.name, func(t *testing.T) {
			handler := newAssetHandlerCreator(assetUseCase).Scan()
			if err := handler.Do(tt.args.action, tt.args.dispatch); (err != nil) != tt.wantErr {
				t.Errorf("Do() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
