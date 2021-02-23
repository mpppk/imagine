package queryimpl

import (
	"encoding/json"
	"fmt"

	"github.com/mpppk/imagine/domain/model"
	"github.com/mpppk/imagine/infra/blt"
	bolt "go.etcd.io/bbolt"
)

type BBoltBaseTag struct {
	bucketName string
	base       *blt.Repository
}

func NewBBoltBaseTag(b *bolt.DB, bucketName string) *BBoltBaseTag {
	return &BBoltBaseTag{
		bucketName: bucketName,
		base:       blt.NewRepository(b),
	}
}

// FIXME: duplicated code
func (b *BBoltTag) Init(ws model.WSName) error {
	if err := b.base.CreateBucketIfNotExist(blt.CreateTagBucketNames(ws)); err != nil {
		return fmt.Errorf("failed to create tag bucket: %w", err)
	}
	if err := b.base.CreateBucketIfNotExist(blt.CreateTagHistoryBucketNames(ws)); err != nil {
		return fmt.Errorf("failed to create tag history bucket: %w", err)
	}
	return nil
}

func (b *BBoltBaseTag) createBucketNames(ws model.WSName) []string {
	return []string{string(ws), b.bucketName}
}

func (b *BBoltBaseTag) loBucketFunc(ws model.WSName, f func(bucket *bolt.Bucket) error) error {
	return b.base.LoBucketFunc(b.createBucketNames(ws), f)
}

func (b *BBoltBaseTag) Get(ws model.WSName, id model.TagID) (tagWithIndex *model.TagWithIndex, exist bool, err error) {
	data, exist, err := b.base.Get(b.createBucketNames(ws), uint64(id))
	if err != nil {
		return nil, exist, err
	}
	if !exist {
		return nil, false, nil
	}

	var a model.TagWithIndex
	if err := json.Unmarshal(data, &a); err != nil {
		return nil, exist, fmt.Errorf("failed to unmarshal json to tag. contents: %s: %w", string(data), err)
	}
	return &a, exist, nil
}

func (b *BBoltBaseTag) ListByAsync(ws model.WSName, f func(tagWithIndex *model.TagWithIndex) bool, cap int) (assetChan <-chan *model.TagWithIndex, err error) {
	c := make(chan *model.TagWithIndex, cap)
	ec := make(chan error, 1)
	f2 := f
	if f2 == nil {
		f2 = func(tagWithIndex *model.TagWithIndex) bool {
			return true
		}
	}
	eachF := func(tagWithIndex *model.TagWithIndex) error {
		if f2(tagWithIndex) {
			c <- tagWithIndex
		}
		return nil
	}

	go func() {
		if err := b.ForEach(ws, eachF); err != nil {
			ec <- fmt.Errorf("failed to list assets: %w", err)
		}
		close(c)
		close(ec)
	}()
	return c, nil
}

func (b *BBoltBaseTag) ListAll(ws model.WSName) (assets []*model.TagWithIndex, err error) {
	return b.ListBy(ws, func(tag *model.TagWithIndex) bool { return true })
}

func (b *BBoltBaseTag) ListBy(ws model.WSName, f func(tag *model.TagWithIndex) bool) (assets []*model.TagWithIndex, err error) {
	eachF := func(tagWithIndex *model.TagWithIndex) error {
		if f(tagWithIndex) {
			assets = append(assets, tagWithIndex)
		}
		return nil
	}
	if err := b.ForEach(ws, eachF); err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}
	return
}

func (b *BBoltBaseTag) ListAsSet(ws model.WSName) (set *model.TagSet, err error) {
	tags, err := b.ListAll(ws)
	if err != nil {
		return nil, err
	}

	set = model.NewTagSet(nil)
	for _, tag := range tags {
		set.Set(tag)
	}
	return
}

func (b *BBoltBaseTag) ForEach(ws model.WSName, f func(tagWithIndex *model.TagWithIndex) error) error {
	return b.loBucketFunc(ws, func(bucket *bolt.Bucket) error {
		return bucket.ForEach(func(k, v []byte) error {
			var tagWithIndex model.TagWithIndex
			if err := json.Unmarshal(v, &tagWithIndex); err != nil {
				return fmt.Errorf("failed to unmarshal tag: %w", err)
			}
			return f(&tagWithIndex)
		})
	})
}
