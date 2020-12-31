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

func SplitIfHasID(assets []*model.ImportAsset) (assetsWithID, assetsWithOutID []*model.ImportAsset) {
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
