package repoimpl

import (
	"fmt"

	"github.com/mpppk/imagine/infra/blt"

	"github.com/blang/semver/v4"
	"github.com/mpppk/imagine/domain/repository"
	bolt "go.etcd.io/bbolt"
)

type BoltMeta struct {
	base        *blt.Repository
	bucketNames []string
}

const (
	versionKey = "version"
)

func (b *BoltMeta) GetDBVersion() (v *semver.Version, exist bool, err error) {
	err = b.loBucketFunc(func(bucket *bolt.Bucket) error {
		rawV := bucket.Get([]byte(versionKey))
		if rawV == nil {
			return nil
		}
		version, err := semver.Parse(string(rawV))
		if err != nil {
			return fmt.Errorf("failed to parse version from %s: %w", string(rawV), err)
		}
		v = &version
		exist = true
		return nil
	})
	return
}

func (b *BoltMeta) SetDBVersion(version *semver.Version) error {
	return b.bucketFunc(func(bucket *bolt.Bucket) error {
		if err := bucket.Put([]byte(versionKey), []byte(version.String())); err != nil {
			return fmt.Errorf("failed to put version: %w", err)
		}
		return nil
	})
}

func NewBoltMeta(b *bolt.DB) repository.Meta {
	return &BoltMeta{
		base:        blt.NewRepository(b),
		bucketNames: blt.CreateMetaBucketNames(),
	}
}

func (b *BoltMeta) Init() error {
	return b.base.CreateBucketIfNotExist(blt.CreateMetaBucketNames())
}

func (b *BoltMeta) loBucketFunc(f func(bucket *bolt.Bucket) error) error {
	return b.base.LoBucketFunc(b.bucketNames, f)
}

func (b *BoltMeta) bucketFunc(f func(bucket *bolt.Bucket) error) error {
	return b.base.BucketFunc(b.bucketNames, f)
}
