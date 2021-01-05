//go:generate mockgen -source=./asset.go -destination=../../infra/repoimpl/mock_asset.go -package=repoimpl

package repository

import (
	"context"

	"github.com/mpppk/imagine/domain/model"
)

type Asset interface {
	Init(ws model.WSName) error
	Add(ws model.WSName, asset *model.Asset) (model.AssetID, error)
	BatchAdd(ws model.WSName, assets []*model.Asset) ([]model.AssetID, error)
	Close() error
	AddByFilePathIfDoesNotExist(ws model.WSName, filePath string) (id model.AssetID, added bool, err error)
	AddByFilePathListIfDoesNotExist(ws model.WSName, filePathList []string) (idList []model.AssetID, err error)
	Get(ws model.WSName, id model.AssetID) (asset *model.Asset, exist bool, err error)
	GetByPath(ws model.WSName, path string) (asset *model.Asset, exist bool, err error)
	Has(ws model.WSName, id model.AssetID) (ok bool, err error)
	Update(ws model.WSName, asset *model.Asset) error
	BatchUpdateByID(ws model.WSName, assets []*model.Asset) (updatedAssets, skippedAssets []*model.Asset, err error)
	BatchUpdateByPath(ws model.WSName, assets []*model.Asset) (updatedAssets, skippedAssets []*model.Asset, err error)
	Delete(ws model.WSName, id model.AssetID) error
	ListByAsync(ctx context.Context, ws model.WSName, f func(asset *model.Asset) bool, cap int) (assetChan <-chan *model.Asset, err error)
	ListRawByAsync(ctx context.Context, ws model.WSName, f func(v []byte) bool, cap int) (c <-chan []byte, err error)
	List(ws model.WSName) (assets []*model.Asset, err error)
	ListBy(ws model.WSName, f func(asset *model.Asset) bool) (assets []*model.Asset, err error)
	ListByIDList(ws model.WSName, idList []model.AssetID) ([]*model.Asset, error)
	ListByPaths(ws model.WSName, paths []string) ([]*model.Asset, error)
	ListByIDListAsync(ctx context.Context, ws model.WSName, idList []model.AssetID, cap int) (assetChan <-chan *model.Asset, errChan <-chan error, err error)
	ForEach(ws model.WSName, f func(asset *model.Asset) error) error
	Revalidate(ws model.WSName) error
}
