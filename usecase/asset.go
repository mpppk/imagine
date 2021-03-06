//go:generate mockgen -source=./asset.go -destination=./mock_usecase/mock_asset.go

package usecase

import (
	"context"
	"io"

	"github.com/mpppk/imagine/domain/model"
)

type Asset interface {
	Init(ws model.WSName) error
	AddAssetFromImagePathListIfDoesNotExist(ws model.WSName, filePathList []string) ([]model.AssetID, error)
	ReadImportAssetsWithProgressBar(ws model.WSName, reader io.Reader, capacity int, f func(assets []*model.ImportAsset) error) error
	SaveImportAssetsFromReader(ws model.WSName, reader io.Reader, capacity int, queries []*model.Query) error
	SaveImportAssets(ws model.WSName, importAssets []*model.ImportAsset, queries []*model.Query) (addedAssets, updatedAssets, skippedAssets []*model.Asset, err error)
	BatchUpdateByID(ws model.WSName, assets []*model.Asset, queries []*model.Query) (updatedAssets, filteredAssets, skippedAssets []*model.Asset, err error)
	BatchUpdateByPath(ws model.WSName, assets []*model.Asset, queries []*model.Query) (updatedAssets, filteredAssets, skippedAssets []*model.Asset, err error)
	AssignBoundingBox(ws model.WSName, assetID model.AssetID, box *model.BoundingBox) (*model.Asset, error)
	UnAssignBoundingBox(ws model.WSName, assetID model.AssetID, boxID model.BoundingBoxID) (*model.Asset, error)
	ModifyBoundingBox(ws model.WSName, assetID model.AssetID, box *model.BoundingBox) (*model.Asset, error)
	DeleteBoundingBox(ws model.WSName, assetID model.AssetID, boxID model.BoundingBoxID) error
	ListAsyncByQueries(ctx context.Context, ws model.WSName, queries []*model.Query) (<-chan *model.Asset, error)
	ListAsyncWithFormat(wsName model.WSName, formatType string, capacity int) (<-chan string, <-chan error, error)
}
