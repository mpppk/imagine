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
	Bolt *bolt.DB
}

func NewBBoltAsset(b *bolt.DB) repository.Asset {
	return &BBoltAsset{
		Bolt: b,
	}
}

func (b *BBoltAsset) Init() error {
	return b.createBucketIfNotExist()
}

func (b *BBoltAsset) createBucketIfNotExist() error {
	return b.Bolt.Update(func(tx *bolt.Tx) error {
		if bucket := tx.Bucket([]byte(assetBucketName)); bucket != nil {
			return nil
		}
		_, err := tx.CreateBucket([]byte(assetBucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}

func (b *BBoltAsset) Close() error {
	return b.Bolt.Close()
}

func (b *BBoltAsset) Add(asset *model.Asset) error {
	return b.bucketFunc(func(bucket *bolt.Bucket) error {
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

func (b *BBoltAsset) Get(id model.AssetID) (asset *model.Asset, err error) {
	err = b.loBucketFunc(func(bucket *bolt.Bucket) error {
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

func (b *BBoltAsset) Update(asset *model.Asset) error {
	return b.bucketFunc(func(bucket *bolt.Bucket) error {
		s, err := json.Marshal(asset)
		if err != nil {
			return fmt.Errorf("failed to marshal asset to json: %w", err)
		}
		return bucket.Put(b.itob(asset.ID), s)
	})
}

func (b *BBoltAsset) ListByAsync(f func(asset *model.Asset) bool, cap int) (assetChan <-chan *model.Asset, err error) {
	c := make(chan *model.Asset, cap)
	eachF := func(asset *model.Asset) error {
		if f(asset) {
			c <- asset
		}
		return nil
	}
	if err := b.ForEach(eachF); err != nil {
		return nil, fmt.Errorf("failed to list assets: %w", err)
	}
	return c, nil
}

func (b *BBoltAsset) ListBy(f func(asset *model.Asset) bool) (assets []*model.Asset, err error) {
	eachF := func(asset *model.Asset) error {
		if f(asset) {
			assets = append(assets, asset)
		}
		return nil
	}
	if err := b.ForEach(eachF); err != nil {
		return nil, fmt.Errorf("failed to list assets: %w", err)
	}
	return
}

func (b *BBoltAsset) ListByTags(tags []model.Tag) (assets []*model.Asset, err error) {
	if len(tags) == 0 {
		return nil, errors.New("no tags given to ListByTags")
	}
	return b.ListBy(func(asset *model.Asset) bool {
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

func (b *BBoltAsset) ForEach(f func(asset *model.Asset) error) error {
	return b.loBucketFunc(func(bucket *bolt.Bucket) error {
		return bucket.ForEach(func(k, v []byte) error {
			var asset model.Asset
			if err := json.Unmarshal(v, &asset); err != nil {
				return fmt.Errorf("failed to unmarshal asset: %w", err)
			}
			return f(&asset)
		})
	})
}

func (b *BBoltAsset) bucket(tx *bolt.Tx) *bolt.Bucket {
	return tx.Bucket([]byte(assetBucketName))
}

func (b *BBoltAsset) bucketFunc(f func(bucket *bolt.Bucket) error) error {
	return b.Bolt.Update(func(tx *bolt.Tx) error {
		return f(b.bucket(tx))
	})
}

func (b *BBoltAsset) loBucketFunc(f func(bucket *bolt.Bucket) error) error {
	return b.Bolt.View(func(tx *bolt.Tx) error {
		return f(b.bucket(tx))
	})
}

func (b *BBoltAsset) itob(id model.AssetID) []byte {
	return itob(uint64(id))
}
