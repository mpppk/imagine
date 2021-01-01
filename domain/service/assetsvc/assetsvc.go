package assetsvc

import (
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
		if asset.ID == 0 {
			assetsWithOutID = append(assetsWithOutID, asset)
		} else {
			assetsWithID = append(assetsWithID, asset)
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

func ToPaths(assets []*model.Asset) (paths []string) {
	for _, asset := range assets {
		paths = append(paths, asset.Path)
	}
	return
}
