package repository

import (
	"github.com/mpppk/imagine/domain/model"
)

type Asset interface {
	Init(ws model.WSName) error
	Close(ws model.WSName) error
	Add(ws model.WSName, asset *model.Asset) error
	Get(ws model.WSName, id model.AssetID) (asset *model.Asset, err error)
	Update(ws model.WSName, asset *model.Asset) error
	ListByAsync(ws model.WSName, f func(asset *model.Asset) bool, cap int) (assetChan <-chan *model.Asset, err error)
	ListBy(ws model.WSName, f func(asset *model.Asset) bool) (assets []*model.Asset, err error)
	ListByTags(ws model.WSName, tags []model.Tag) (assets []*model.Asset, err error)
	ForEach(ws model.WSName, f func(asset *model.Asset) error) error
}
