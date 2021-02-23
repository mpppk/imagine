package queryimpl

import (
	"github.com/mpppk/imagine/domain/query"
	"github.com/mpppk/imagine/infra/blt"
	bolt "go.etcd.io/bbolt"
)

type BBoltTag struct {
	*BBoltBaseTag
}

func NewBBoltTag(b *bolt.DB) query.Tag {
	return &BBoltTag{
		BBoltBaseTag: NewBBoltBaseTag(b, blt.TagBucketName),
	}
}
