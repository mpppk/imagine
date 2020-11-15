package repository

import "github.com/blang/semver/v4"

type Meta interface {
	Init() error
	GetDBVersion() (*semver.Version, bool, error)
	SetDBVersion(version *semver.Version) error
}
