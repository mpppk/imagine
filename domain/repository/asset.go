package repository

import (
	"github.com/mpppk/imagine/domain/model"
)

type Asset interface {
	Init() error
	Close() error
	Add(*model.Asset) error
	Get(id model.AssetID) (asset *model.Asset, err error)
	Update(asset *model.Asset) error
	ListByAsync(f func(asset *model.Asset) bool, cap int) (assetChan <-chan *model.Asset, err error)
	ListBy(f func(asset *model.Asset) bool) (assets []*model.Asset, err error)
	ListByTags(tags []model.Tag) (assets []*model.Asset, err error)
	ForEach(f func(asset *model.Asset) error) error
}
