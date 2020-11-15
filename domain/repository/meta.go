package repository

import "github.com/blang/semver/v4"

type Meta interface {
	Init() error
	GetVersion() (*semver.Version, bool, error)
	SetVersion(version *semver.Version) error
	//CompareVersion() (c int, appV, dbV *semver.Version, err error)
}
