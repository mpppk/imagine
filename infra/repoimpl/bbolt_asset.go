package repoimpl

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mpppk/imagine/domain/repository"

	"github.com/mpppk/imagine/domain/model"
	bolt "go.etcd.io/bbolt"
)

const assetBucketName = "Asset"

type BBoltAsset struct {
	base           *boltRepository
	pathRepository *pathRepository
}

func createBucketNames(ws model.WSName) []string {
	return []string{string(ws), assetBucketName}
}

func NewBBoltAsset(b *bolt.DB) repository.Asset {
	return &BBoltAsset{
		base:           newBoltRepository(b),
		pathRepository: newPathRepository(b),
	}
}

func (b *BBoltAsset) Init(ws model.WSName) error {
	return nil
	//return b.base.createBucketIfNotExist(createBucketNames(ws))
}

func (b *BBoltAsset) Close(ws model.WSName) error {
	return b.base.close()
}

func (b *BBoltAsset) Add(ws model.WSName, asset *model.Asset) error {
	return b.base.bucketFunc(createBucketNames(ws), func(bucket *bolt.Bucket) error {
		id, err := bucket.NextSequence()
		if err != nil {
			return err
		}
		asset.ID = model.AssetID(id)
		s, err := json.Marshal(asset)
		if err != nil {
			return fmt.Errorf("failed to marshal asset to json: %w", err)
		}
		return bucket.Put(b.itob(asset.ID), s)
	})
}

func (b *BBoltAsset) Get(ws model.WSName, id model.AssetID) (asset *model.Asset, err error) {
	// FIXME: convert to read only
	err = b.base.bucketFunc(createBucketNames(ws), func(bucket *bolt.Bucket) error {
		data := bucket.Get(b.itob(id))
		if data == nil {
			return fmt.Errorf("a does not exist: %v", id)
		}
		var a model.Asset
		if err := json.Unmarshal(data, &a); err != nil {
			return err
		}
		asset = &a
		return nil
	})
	return
}

func (b *BBoltAsset) Update(ws model.WSName, asset *model.Asset) error {
	return b.base.bucketFunc(createBucketNames(ws), func(bucket *bolt.Bucket) error {
		s, err := json.Marshal(asset)
		if err != nil {
			return fmt.Errorf("failed to marshal asset to json: %w", err)
		}
		return bucket.Put(b.itob(asset.ID), s)
	})
}

func (b *BBoltAsset) ListByAsync(ws model.WSName, f func(asset *model.Asset) bool, cap int) (assetChan <-chan *model.Asset, err error) {
	c := make(chan *model.Asset, cap)
	ec := make(chan error, 1)
	f2 := f
	if f2 == nil {
		f2 = func(asset *model.Asset) bool {
			return true
		}
	}
	eachF := func(asset *model.Asset) error {
		if f2(asset) {
			c <- asset
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

func (b *BBoltAsset) ListBy(ws model.WSName, f func(asset *model.Asset) bool) (assets []*model.Asset, err error) {
	eachF := func(asset *model.Asset) error {
		if f(asset) {
			assets = append(assets, asset)
		}
		return nil
	}
	if err := b.ForEach(ws, eachF); err != nil {
		return nil, fmt.Errorf("failed to list assets: %w", err)
	}
	return
}

func (b *BBoltAsset) ListByTags(ws model.WSName, tags []model.Tag) (assets []*model.Asset, err error) {
	if len(tags) == 0 {
		return nil, errors.New("no tags given to ListByTags")
	}
	return b.ListBy(ws, func(asset *model.Asset) bool {
		m := map[model.Tag]struct{}{}
		for _, t := range asset.Tags {
			m[t] = struct{}{}
		}
		for _, tag := range tags {
			if _, ok := m[tag]; !ok {
				return false
			}
		}
		return true
	})
}

func (b *BBoltAsset) ForEach(ws model.WSName, f func(asset *model.Asset) error) error {
	// FIXME: convert to read only
	return b.base.bucketFunc(createBucketNames(ws), func(bucket *bolt.Bucket) error {
		return bucket.ForEach(func(k, v []byte) error {
			var asset model.Asset
			if err := json.Unmarshal(v, &asset); err != nil {
				return fmt.Errorf("failed to unmarshal asset: %w", err)
			}
			return f(&asset)
		})
	})
}

func (b *BBoltAsset) itob(id model.AssetID) []byte {
	return itob(uint64(id))
}
