package repoimpl

import (
	"encoding/json"
	"fmt"

	"github.com/mpppk/imagine/domain/repository"

	"github.com/mpppk/imagine/domain/model"
	bolt "go.etcd.io/bbolt"
)

type BBoltTag struct {
	base *boltRepository
}

func NewBBoltTag(b *bolt.DB) repository.Tag {
	return &BBoltTag{
		base: newBoltRepository(b),
	}
}

func (b *BBoltTag) loBucketFunc(ws model.WSName, f func(bucket *bolt.Bucket) error) error {
	return b.base.loBucketFunc(createTagBucketNames(ws), f)
}

func (b *BBoltTag) Init(ws model.WSName) error {
	return b.base.createBucketIfNotExist(createTagBucketNames(ws))
}

func (b *BBoltTag) Add(ws model.WSName, tag *model.Tag) error {
	return b.base.add(createTagBucketNames(ws), tag)
}

func (b *BBoltTag) Get(ws model.WSName, id model.TagID) (tag *model.Tag, err error) {
	data, err := b.base.get(createTagBucketNames(ws), uint64(id))
	if err != nil {
		return nil, err
	}
	var a model.Tag
	if err := json.Unmarshal(data, &a); err != nil {
		return nil, err
	}
	return &a, nil
}

func (b *BBoltTag) RecreateBucket(ws model.WSName) error {
	return b.base.recreateBucket(createTagBucketNames(ws))
}

func (b *BBoltTag) Update(ws model.WSName, tag *model.Tag) error {
	return b.base.update(createTagBucketNames(ws), tag)
}

func (b *BBoltTag) ListByAsync(ws model.WSName, f func(tag *model.Tag) bool, cap int) (assetChan <-chan *model.Tag, err error) {
	c := make(chan *model.Tag, cap)
	ec := make(chan error, 1)
	f2 := f
	if f2 == nil {
		f2 = func(tag *model.Tag) bool {
			return true
		}
	}
	eachF := func(tag *model.Tag) error {
		if f2(tag) {
			c <- tag
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

func (b *BBoltTag) ListAll(ws model.WSName) (assets []*model.Tag, err error) {
	return b.ListBy(ws, func(tag *model.Tag) bool { return true })
}

func (b *BBoltTag) ListBy(ws model.WSName, f func(tag *model.Tag) bool) (assets []*model.Tag, err error) {
	eachF := func(tag *model.Tag) error {
		if f(tag) {
			assets = append(assets, tag)
		}
		return nil
	}
	if err := b.ForEach(ws, eachF); err != nil {
		return nil, fmt.Errorf("failed to list assets: %w", err)
	}
	return
}

func (b *BBoltTag) ForEach(ws model.WSName, f func(tag *model.Tag) error) error {
	return b.loBucketFunc(ws, func(bucket *bolt.Bucket) error {
		return bucket.ForEach(func(k, v []byte) error {
			var tag model.Tag
			if err := json.Unmarshal(v, &tag); err != nil {
				return fmt.Errorf("failed to unmarshal tag: %w", err)
			}
			return f(&tag)
		})
	})
}
