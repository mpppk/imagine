package repoimpl

import (
	"bytes"
	"context"
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

func (b *BBoltAsset) AddByFilePathListIfDoesNotExist(ws model.WSName, filePathList []string) ([]model.AssetID, error) {
	notExistPaths, err := b.pathRepository.FilterExistPath(ws, filePathList)
	if err != nil {
		return nil, err
	}

	var assets []*model.Asset
	for _, p := range notExistPaths {
		assets = append(assets, model.NewAssetFromFilePath(p))
	}

	idList, err := b.AddList(ws, assets)
	if err != nil {
		return nil, err
	}

	return idList, b.pathRepository.AddList(ws, notExistPaths, idList)
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
	return id, true, nil
}

func (b *BBoltAsset) AddList(ws model.WSName, assets []*model.Asset) ([]model.AssetID, error) {
	var dataList []boltData
	for _, asset := range assets {
		dataList = append(dataList, asset)
	}
	idList, err := b.base.addListByID(createAssetBucketNames(ws), dataList)
	if err != nil {
		return nil, err
	}

	var assetIDList []model.AssetID
	for _, id := range idList {
		assetIDList = append(assetIDList, model.AssetID(id))
	}

	return assetIDList, nil
}

func (b *BBoltAsset) Add(ws model.WSName, asset *model.Asset) (model.AssetID, error) {
	id, err := b.base.addByID(createAssetBucketNames(ws), asset)
	if err != nil {
		return 0, err
	}
	return model.AssetID(id), b.pathRepository.Add(ws, asset.Path, model.AssetID(id))
}

func (b *BBoltAsset) Get(ws model.WSName, id model.AssetID) (asset *model.Asset, exist bool, err error) {
	data, exist, err := b.base.get(createAssetBucketNames(ws), uint64(id))
	if err != nil {
		return nil, false, err
	}
	var a model.Asset
	if err := json.Unmarshal(data, &a); err != nil {
		return nil, false, err
	}
	return &a, exist, nil
}

func (b *BBoltAsset) GetByPath(ws model.WSName, path string) (asset *model.Asset, exist bool, err error) {
	id, exist, err := b.pathRepository.Get(ws, path)
	if err != nil {
		return nil, false, err
	}
	if !exist {
		return nil, false, nil
	}
	return b.Get(ws, id)
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

func (b *BBoltAsset) List(ctx context.Context, ws model.WSName, cap int) (assetChan <-chan *model.Asset, err error) {
	f := func(asset *model.Asset) bool {
		return true
	}
	return b.ListByAsync(ctx, ws, f, cap)
}

func (b *BBoltAsset) ListByAsync(ctx context.Context, ws model.WSName, f func(asset *model.Asset) bool, cap int) (assetChan <-chan *model.Asset, err error) {
	c := make(chan *model.Asset, cap)
	ec := make(chan error, 1)
	f2 := f
	if f2 == nil {
		f2 = func(asset *model.Asset) bool {
			return true
		}
	}

	// FIXME: goroutine leak
	go func() {
		batchNum := 50
		min := itob(0)
	L:
		for {
			var assets []*model.Asset
			var lastAsset *model.Asset = nil
			err := b.base.loBucketFunc(createAssetBucketNames(ws), func(bucket *bolt.Bucket) error {
				cursor := bucket.Cursor()
				cnt := 0
				for k, v := cursor.Seek(min); k != nil && cnt < batchNum; k, v = cursor.Next() {
					if bytes.Equal(k, min) {
						continue
					}
					cnt++
					var asset model.Asset
					if err := json.Unmarshal(v, &asset); err != nil {
						return fmt.Errorf("failed to unmarshal asset: %w", err)
					}
					if f2(&asset) {
						assets = append(assets, &asset)
					}
					lastAsset = &asset
				}
				return nil
			})

			if err != nil {
				ec <- fmt.Errorf("failed to list assets: %w", err)
			}

			if lastAsset == nil {
				break
			}

			for _, asset := range assets {
				select {
				case <-ctx.Done():
					break L
				case c <- asset:
				}
			}
			min = itob(uint64(lastAsset.ID))
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

func (b *BBoltAsset) Revalidate(ws model.WSName) error {
	cap := 10000
	if err := b.pathRepository.DeleteAll(ws); err != nil {
		return fmt.Errorf("failed to delete path caches while revalidating: %w", err)
	}
	c, err := b.List(context.Background(), ws, cap)
	if err != nil {
		return fmt.Errorf("failed to prepare asset listing: %w", err)
	}

	paths := make([]string, 0, cap)
	idList := make([]model.AssetID, 0, cap)
	for asset := range c {
		paths = append(paths, asset.Path)
		idList = append(idList, asset.ID)
		if len(paths) >= cap {
			if err := b.pathRepository.AddList(ws, paths, idList); err != nil {
				return fmt.Errorf("failed to add paths: %w", err)
			}
			paths = make([]string, 0, cap)
			idList = make([]model.AssetID, 0, cap)
		}
	}
	return nil
}
