package repoimpl_test

import (
	"testing"

	"github.com/mpppk/imagine/testutil"

	"github.com/blang/semver/v4"
	"github.com/mpppk/imagine/usecase/usecasetest"
)

func TestBoltMeta_SetAndGetVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		wantErr bool
	}{
		{
			version: "0.0.1",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		usecasetest.RunParallelWithUseCases(t, tt.name, "", func(t *testing.T, ut *usecasetest.UseCases) {
			v, err := semver.New(tt.version)
			if err != nil {
				t.Errorf("failed to create semver struct: %v", err)
			}
			if err := ut.Usecases.Client.Meta.SetDBVersion(v); (err != nil) != tt.wantErr {
				t.Errorf("SetDBVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
			gotVersion, ok, err := ut.Usecases.Client.Meta.GetDBVersion()
			if err != nil || !ok {
				t.Errorf("failed to get version: %v", err)
			}

			testutil.Diff(t, v, gotVersion)
		})
	}
}
