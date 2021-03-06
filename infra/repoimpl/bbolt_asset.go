package repoimpl

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mpppk/imagine/infra/blt"

	"github.com/mpppk/imagine/domain/service/assetsvc"

	"github.com/mpppk/imagine/domain/repository"

	"github.com/mpppk/imagine/domain/model"
	bolt "go.etcd.io/bbolt"
)

type BBoltAsset struct {
	base           *blt.Repository
	pathRepository *bboltPathRepository
}

func NewBBoltAsset(b *bolt.DB) repository.Asset {
	return &BBoltAsset{
		base:           blt.NewRepository(b),
		pathRepository: newBBoltPathRepository(b),
	}
}

func (b *BBoltAsset) Init(ws model.WSName) error {
	if err := b.base.CreateBucketIfNotExist(blt.CreateAssetBucketNames(ws)); err != nil {
		return fmt.Errorf("failed to create asset bucket: %w", err)
	}
	if err := b.base.CreateBucketIfNotExist(blt.CreatePathBucketNames(ws)); err != nil {
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
		asset := model.NewAssetFromFilePath(p)
		assets = append(assets, asset)
	}

	idList, err := b.BatchAdd(ws, assets)
	if err != nil {
		return nil, err
	}

	return idList, nil
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

// BatchAdd add Assets.
// provided assets must not have ID and must have path. If asset does not satisfy them, error will be returned.
func (b *BBoltAsset) BatchAdd(ws model.WSName, assets []*model.Asset) ([]model.AssetID, error) {
	var dataList []blt.BoltData
	var paths []string
	for _, asset := range assets {
		if ok, reason := asset.IsAddable(); !ok {
			return nil, fmt.Errorf("failed to add asset because it is not addable. reason: %v. asset:%#v", reason, asset)
		}

		dataList = append(dataList, asset)
		paths = append(paths, asset.Path)
	}
	idList, err := b.base.BatchAddWithID(blt.CreateAssetBucketNames(ws), dataList)
	if err != nil {
		return nil, err
	}

	assetIDList := toAssetIDList(idList)
	if err := b.pathRepository.AddList(ws, paths, assetIDList); err != nil {
		return nil, fmt.Errorf("failed to add paths to path repository: %w", err)
	}

	return assetIDList, nil
}

func (b *BBoltAsset) Add(ws model.WSName, asset *model.Asset) (model.AssetID, error) {
	id, err := b.base.AddWithID(blt.CreateAssetBucketNames(ws), asset)
	if err != nil {
		return 0, err
	}
	return model.AssetID(id), b.pathRepository.Add(ws, asset.Path, model.AssetID(id))
}

func (b *BBoltAsset) Get(ws model.WSName, id model.AssetID) (asset *model.Asset, exist bool, err error) {
	data, exist, err := b.base.Get(blt.CreateAssetBucketNames(ws), uint64(id))
	if err != nil {
		return nil, false, fmt.Errorf("failed to get asset by ID(%v): %w", id, err)
	}
	if !exist {
		return nil, false, nil
	}

	var a model.Asset
	if err := json.Unmarshal(data, &a); err != nil {
		return nil, false, fmt.Errorf("failed to unmarshal asset. raw text: %s: %w", string(data), err)
	}
	return &a, exist, nil
}

func (b *BBoltAsset) GetByPath(ws model.WSName, path string) (asset *model.Asset, exist bool, err error) {
	id, exist, err := b.pathRepository.Get(ws, path)
	if err != nil {
		return nil, false, fmt.Errorf("failed to get asset from path repository: %w", err)
	}
	if !exist {
		return nil, false, nil
	}
	return b.Get(ws, id)
}

func (b *BBoltAsset) Has(ws model.WSName, id model.AssetID) (ok bool, err error) {
	_, exist, err := b.base.Get(blt.CreateAssetBucketNames(ws), uint64(id))
	return exist, err
}

// Update updates asset by ID.
// if data which have ID does not exist, return error.
func (b *BBoltAsset) Update(ws model.WSName, asset *model.Asset) error {
	return b.base.UpdateByID(blt.CreateAssetBucketNames(ws), asset)
}

// BatchUpdate update assets by ID.
// Invalid asset will be skip. For example, an asset that contains a bounding box that does not have an ID.
func (b *BBoltAsset) BatchUpdateByID(ws model.WSName, assets []*model.Asset) (updatedAssets, skippedAssets []*model.Asset, err error) {
	var dataList []blt.BoltData
	for _, asset := range assets {
		if !asset.IsUpdatableByID() {
			skippedAssets = append(skippedAssets, asset)
			continue
		}
		dataList = append(dataList, asset)
	}
	updatedDataList, skippedDataList, err := b.base.BatchUpdateByID(blt.CreateAssetBucketNames(ws), dataList)
	for _, data := range updatedDataList {
		asset := data.(*model.Asset)
		updatedAssets = append(updatedAssets, asset)
	}
	for _, data := range skippedDataList {
		asset := data.(*model.Asset)
		skippedAssets = append(skippedAssets, asset)
	}
	return
}

// BatchUpdateByPath update assets by path.
// Invalid asset will be skip. For example, an asset that contains a bounding box that does not have an ID.
// If asset which have non exist path is provided, it will be ignored.
func (b *BBoltAsset) BatchUpdateByPath(ws model.WSName, assets []*model.Asset) (updatedAssets, skippedAssets []*model.Asset, err error) {
	assetIDList, err := b.pathRepository.ListByPath(ws, assetsvc.ToPaths(assets))
	if err != nil {
		return nil, nil, err
	}

	for i, asset := range assets {
		asset.ID = assetIDList[i]
	}

	return b.BatchUpdateByID(ws, assets)
}

func (b *BBoltAsset) Delete(ws model.WSName, id model.AssetID) error {
	return b.base.Delete(blt.CreateAssetBucketNames(ws), uint64(id))
}

// ListByIDList list assets by provided ID ID list.
// If ID which does not exist is provided, nil will be returned.
func (b *BBoltAsset) ListByIDList(ws model.WSName, idList []model.AssetID) (assets []*model.Asset, err error) {
	var rawIdList []uint64
	for _, id := range idList {
		rawIdList = append(rawIdList, uint64(id))
	}
	contents, err := b.base.BatchGet(blt.CreateAssetBucketNames(ws), rawIdList)
	if err != nil {
		return nil, err
	}

	for _, content := range contents {
		if content == nil {
			assets = append(assets, nil)
			continue
		}
		asset, err := model.NewAssetFromJson(content)
		if err != nil {
			return nil, fmt.Errorf("failed to create new asset from json: %w", err)
		}
		assets = append(assets, asset)
	}

	return
}

func (b *BBoltAsset) ListByPaths(ws model.WSName, paths []string) (assets []*model.Asset, err error) {
	idList, err := b.pathRepository.ListByPath(ws, paths)
	if err != nil {
		return nil, fmt.Errorf("failed to list asset ID from path repository: %w", err)
	}
	assets, err = b.ListByIDList(ws, idList)
	if err != nil {
		return nil, fmt.Errorf("failed to list assets from asset ID: %w", err)
	}
	return
}

func (b *BBoltAsset) ListAsync(ctx context.Context, ws model.WSName, cap int) (assetChan <-chan *model.Asset, err error) {
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

	go func() {
		batchNum := 50
		min := blt.Itob(0)
	L:
		for {
			var assets []*model.Asset
			var lastAsset *model.Asset = nil
			err := b.base.LoBucketFunc(blt.CreateAssetBucketNames(ws), func(bucket *bolt.Bucket) error {
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
			min = blt.Itob(uint64(lastAsset.ID))
		}
		close(c)
		close(ec)
	}()
	return c, nil
}

func getAssetByIdFromBucket(bucket *bolt.Bucket, id model.AssetID) (*model.Asset, error) {
	v := bucket.Get(blt.Itob(uint64(id)))
	var asset model.Asset
	if err := json.Unmarshal(v, &asset); err != nil {
		return nil, fmt.Errorf("failed to unmarshal asset: %w", err)
	}
	return &asset, nil
}

func (b *BBoltAsset) ListByIDListAsync(ctx context.Context, ws model.WSName, idList []model.AssetID, cap int) (assetChan <-chan *model.Asset, errChan <-chan error, err error) {
	c := make(chan *model.Asset, cap)
	ec := make(chan error, 1)
	go func() {
		index := 0
	L:
		for index < len(idList) {
			var assets []*model.Asset
			err := b.base.LoBucketFunc(blt.CreateAssetBucketNames(ws), func(bucket *bolt.Bucket) error {
				batchCnt := 0
				for batchCnt <= cap {
					if index+batchCnt >= len(idList) {
						break
					}
					asset, err := getAssetByIdFromBucket(bucket, idList[index+batchCnt])
					if err != nil {
						return err
					}
					assets = append(assets, asset)
					batchCnt++
				}
				index += batchCnt
				return nil
			})

			if err != nil {
				ec <- fmt.Errorf("failed to list assets: %w", err)
			}

			for _, asset := range assets {
				select {
				case <-ctx.Done():
					break L
				case c <- asset:
				}
			}
		}
		close(c)
	}()
	return c, ec, nil
}

func (b *BBoltAsset) ListRawByAsync(ctx context.Context, ws model.WSName, f func(v []byte) bool, cap int) (vc <-chan []byte, err error) {
	c := make(chan []byte, cap)
	ec := make(chan error, 1)
	f2 := f
	if f2 == nil {
		f2 = func(v []byte) bool {
			return true
		}
	}

	go func() {
		batchNum := 50
		min := blt.Itob(0)
	L:
		for {
			var assets [][]byte
			var lastAssetID uint64 = 0
			err := b.base.LoBucketFunc(blt.CreateAssetBucketNames(ws), func(bucket *bolt.Bucket) error {
				cursor := bucket.Cursor()
				cnt := 0
				for k, v := cursor.Seek(min); k != nil && cnt < batchNum; k, v = cursor.Next() {
					if bytes.Equal(k, min) {
						continue
					}
					cnt++
					if f2(v) {
						newV := make([]byte, len(v))
						copy(newV, v)
						assets = append(assets, newV)
					}
					lastAssetID = blt.Btoi(k)
				}
				return nil
			})

			if err != nil {
				ec <- fmt.Errorf("failed to list assets: %w", err)
			}

			if lastAssetID == 0 {
				break
			}

			for _, asset := range assets {
				select {
				case <-ctx.Done():
					break L
				case c <- asset:
				}
			}
			min = blt.Itob(lastAssetID)
		}
		close(c)
		close(ec)
	}()
	return c, nil
}

func (b *BBoltAsset) List(ws model.WSName) (assets []*model.Asset, err error) {
	return b.ListBy(ws, func(a *model.Asset) bool { return true })
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

func (b *BBoltAsset) ListByTags(ws model.WSName, tags []model.UnindexedTag) (assets []*model.Asset, err error) {
	if len(tags) == 0 {
		return nil, errors.New("no tags given to ListByTags")
	}
	return b.ListBy(ws, func(asset *model.Asset) bool {
		m := map[model.TagID]struct{}{}
		for _, box := range asset.BoundingBoxes {
			m[box.TagID] = struct{}{}
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
	return b.base.LoBucketFunc(blt.CreateAssetBucketNames(ws), func(bucket *bolt.Bucket) error {
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
	c, err := b.ListAsync(context.Background(), ws, cap)
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

func (b *BBoltAsset) Close() error {
	return b.base.Close()
}

func toAssetIDList(idList []uint64) (assetIDList []model.AssetID) {
	for _, id := range idList {
		assetIDList = append(assetIDList, model.AssetID(id))
	}
	return
}

// UnAssignTag unassign given tag from all assets.
// return assets which have given tag.
func (b *BBoltAsset) UnAssignTags(ws model.WSName, tagIDList ...model.TagID) error {
	if len(tagIDList) == 0 {
		return nil
	}
	f := func(asset *model.Asset) *model.Asset {
		dirty := false
		for _, tagID := range tagIDList {
			if unAssigned := asset.UnAssignTagIfExist(tagID); unAssigned {
				dirty = true
			}
		}
		if !dirty {
			return nil
		}
		return asset
	}
	return b.Map(ws, f)
}

// Map updates each asset by provided function.
// If provided function returns nil, the asset will not be updated.
func (b *BBoltAsset) Map(ws model.WSName, f func(asset *model.Asset) *model.Asset) error {
	return b.base.BucketFunc(blt.CreateAssetBucketNames(ws), func(bucket *bolt.Bucket) error {
		cursor := bucket.Cursor()
		min := blt.Itob(1)
		for k, v := cursor.Seek(min); k != nil; k, v = cursor.Next() {
			asset, err := model.NewAssetFromJson(v)
			if err != nil {
				return fmt.Errorf("failed to unmarshal asset: %w", err)
			}
			if newAsset := f(asset); newAsset != nil {
				b, err := newAsset.ToJsonBytes()
				if err != nil {
					return err
				}
				if err := bucket.Put(k, b); err != nil {
					return fmt.Errorf("failed to map asset: %w", err)
				}
			}
		}
		return nil
	})
}
