package assetsvc

import (
	"fmt"

	"github.com/mpppk/imagine/domain/model"
	"github.com/mpppk/imagine/domain/service/boxsvc"
)

func ToUniqTagNames(assets []*model.ImportAsset) (tagNames []string) {
	m := map[string]struct{}{}
	for _, asset := range assets {
		for _, tagName := range boxsvc.ToUniqTagNames(asset.BoundingBoxes) {
			m[tagName] = struct{}{}
		}
	}
	for tagName := range m {
		tagNames = append(tagNames, tagName)
	}
	return
}

func SplitBy(assets []*model.Asset, f func(asset *model.Asset) bool) (trueAssets, falseAssets []*model.Asset) {
	for _, asset := range assets {
		if f(asset) {
			trueAssets = append(trueAssets, asset)
		} else {
			falseAssets = append(falseAssets, asset)
		}
	}
	return
}

func SplitByPath(assets []*model.Asset) (assetsWithPath, assetsWithOutPath []*model.Asset) {
	return SplitBy(assets, func(asset *model.Asset) bool {
		return asset.HasPath()
	})
}

func SplitByID(assets []*model.Asset) (assetsWithID, assetsWithOutID []*model.Asset) {
	for _, asset := range assets {
		if asset.HasID() {
			assetsWithID = append(assetsWithID, asset)
		} else {
			assetsWithOutID = append(assetsWithOutID, asset)
		}
	}
	return
}

func ToAssets(importAssets []*model.ImportAsset, tagSet *model.TagSet) (assets []*model.Asset, err error) {
	for _, importAsset := range importAssets {
		asset, err := importAsset.ToAsset(tagSet)
		if err != nil {
			return nil, err
		}
		assets = append(assets, asset)
	}
	return
}

func ToAssetIDList(assets []*model.Asset) (assetIDList []model.AssetID) {
	for _, asset := range assets {
		assetIDList = append(assetIDList, asset.ID)
	}
	return
}

// Update merge provided base assets and other assets.
func Update(baseAssets, otherAssets []*model.Asset) error {
	if len(baseAssets) != len(otherAssets) {
		return fmt.Errorf("invalid arguments are given to assetsvc.UpdateBy. length are different(%d, %d)", len(baseAssets), len(otherAssets))
	}
	for i, baseAsset := range baseAssets {
		if err := baseAsset.UpdateBy(otherAssets[i]); err != nil {
			return err
		}
	}
	return nil
}

func ToPaths(assets []*model.Asset) (paths []string) {
	for _, asset := range assets {
		paths = append(paths, asset.Path)
	}
	return
}

func FilterNil(assets []*model.Asset) (filteredAssets []*model.Asset) {
	for _, asset := range assets {
		if asset != nil {
			filteredAssets = append(filteredAssets, asset)
		}
	}
	return
}

// NilIndex return indices which position of nil no provided assets.
func NilIndex(assets []*model.Asset) (indices []int) {
	for i, asset := range assets {
		if asset == nil {
			indices = append(indices, i)
		}
	}
	return
}

func FilterByIndex(assets []*model.Asset, indices []int) (newAssets []*model.Asset) {
	m := map[int]struct{}{}
	for _, index := range indices {
		m[index] = struct{}{}
	}

	for i, asset := range assets {
		if _, ok := m[i]; ok {
			newAssets = append(newAssets, asset)
		}
	}
	return
}

func matchAllQueries(asset *model.Asset, tagSet *model.TagSet, queries []*model.Query, matchToNil bool) bool {
	if asset == nil {
		return matchToNil
	}
	for _, query := range queries {
		if !query.Match(asset, tagSet) {
			return false
		}
	}
	return true
}

// Query filter assets by queries and return matched assets.
// if asset is nil, it will be filtered.
func Query(assets []*model.Asset, queries []*model.Query, tagSet *model.TagSet, matchToNil bool) (matchedAssets, filteredAssets []*model.Asset) {
	for _, asset := range assets {
		if matchAllQueries(asset, tagSet, queries, matchToNil) {
			matchedAssets = append(matchedAssets, asset)
		} else {
			filteredAssets = append(filteredAssets, asset)
		}
	}
	return
}

// QueryIndex filter assets by queries and return matched asset indices.
// if asset is nil, it will be filtered.
func QueryIndex(assets []*model.Asset, queries []*model.Query, tagSet *model.TagSet, matchToNil bool) (matchedAssetsIndex, filteredAssetsIndex []int) {
	for i, asset := range assets {
		if matchAllQueries(asset, tagSet, queries, matchToNil) {
			matchedAssetsIndex = append(matchedAssetsIndex, i)
		} else {
			filteredAssetsIndex = append(filteredAssetsIndex, i)
		}
	}
	return
}
