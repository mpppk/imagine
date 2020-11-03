package repoimpl

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mpppk/imagine/domain/repository"

	"github.com/mpppk/imagine/domain/model"
	bolt "go.etcd.io/bbolt"
)

type BBoltAsset struct {
	base           *boltRepository
	pathRepository *bboltPathRepository
}

func NewBBoltAsset(b *bolt.DB) repository.Asset {
	return &BBoltAsset{
		base:           newBoltRepository(b),
		pathRepository: newBBoltPathRepository(b),
	}
}

func (b *BBoltAsset) Init(ws model.WSName) error {
	if err := b.base.createBucketIfNotExist(createAssetBucketNames(ws)); err != nil {
		return fmt.Errorf("failed to create asset bucket: %w", err)
	}
	if err := b.base.createBucketIfNotExist(createPathBucketNames(ws)); err != nil {
		return fmt.Errorf("failed to create path bucket: %w", err)
	}
	return nil
}

func (b *BBoltAsset) AddByFilePathIfDoesNotExist(ws model.WSName, filePath string) (model.AssetID, bool, error) {
	if _, exist, err := b.pathRepository.Get(ws, filePath); err != nil {
		return 0, false, fmt.Errorf("failed to register asset path: %w", err)
	} else if exist {
		return 0, false, nil
	}

	id, err := b.Add(ws, model.NewAssetFromFilePath(filePath))
	if err != nil {
		return 0, false, err
	}
	return id, true, b.pathRepository.Add(ws, filePath, id)
}

func (b *BBoltAsset) Add(ws model.WSName, asset *model.Asset) (model.AssetID, error) {
	id, err := b.base.addByID(createAssetBucketNames(ws), asset)
	if err != nil {
		return 0, err
	}
	return model.AssetID(id), nil
}

func (b *BBoltAsset) Get(ws model.WSName, id model.AssetID) (asset *model.Asset, err error) {
	data, _, err := b.base.get(createAssetBucketNames(ws), uint64(id))
	if err != nil {
		return nil, err
	}
	var a model.Asset
	if err := json.Unmarshal(data, &a); err != nil {
		return nil, err
	}
	return &a, nil
}

func (b *BBoltAsset) Has(ws model.WSName, id model.AssetID) (ok bool, err error) {
	_, exist, err := b.base.get(createAssetBucketNames(ws), uint64(id))
	return exist, err
}

func (b *BBoltAsset) Update(ws model.WSName, asset *model.Asset) error {
	return b.base.updateByID(createAssetBucketNames(ws), asset)
}

func (b *BBoltAsset) Delete(ws model.WSName, id model.AssetID) error {
	return b.base.delete(createAssetBucketNames(ws), uint64(id))
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
		m := map[model.TagID]struct{}{}
		for _, box := range asset.BoundingBoxes {
			m[box.Tag.ID] = struct{}{}
		}
		for _, tag := range tags {
			if _, ok := m[tag.ID]; !ok {
				return false
			}
		}
		return true
	})
}

func (b *BBoltAsset) ForEach(ws model.WSName, f func(asset *model.Asset) error) error {
	return b.base.loBucketFunc(createAssetBucketNames(ws), func(bucket *bolt.Bucket) error {
		return bucket.ForEach(func(k, v []byte) error {
			var asset model.Asset
			if err := json.Unmarshal(v, &asset); err != nil {
				return fmt.Errorf("failed to unmarshal asset: %w", err)
			}
			return f(&asset)
		})
	})
}
