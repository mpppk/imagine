package repoimpl_test

import (
	"reflect"
	"testing"

	"github.com/blang/semver/v4"
	"github.com/mpppk/imagine/usecase/usecasetest"
)

func TestBoltMeta_SetAndGetVersion(t *testing.T) {
	fileName := "TestBoltMeta_SetAndGetVersion.db"
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
		t.Run(tt.name, func(t *testing.T) {
			usecases, closer, remover := usecasetest.SetUpUseCases(t, fileName, "")
			defer closer()
			defer remover()

			v, err := semver.New(tt.version)
			if err != nil {
				t.Errorf("failed to create semver struct: %v", err)
			}
			if err := usecases.Client.Meta.SetDBVersion(v); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
			gotVersion, ok, err := usecases.Client.Meta.GetDBVersion()
			if err != nil || !ok {
				t.Errorf("failed to get version: %v", err)
			}
			if !reflect.DeepEqual(v, gotVersion) {
				t.Errorf("want: %#v, got: %#v", v, gotVersion)
			}
		})
	}
}
