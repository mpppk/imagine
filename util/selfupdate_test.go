package util

import (
	"reflect"
	"testing"

	semverv3 "github.com/blang/semver"
	"github.com/blang/semver/v4"
)

func Test_toV3PRVersion(t *testing.T) {
	type args struct {
		v semver.Version
	}
	tests := []struct {
		name string
		args args
		want semverv3.Version
	}{
		{
			args: args{semver.MustParse("1.2.3")},
			want: semverv3.MustParse("1.2.3"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toV3PRVersion(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toV3PRVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
