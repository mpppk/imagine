package usecase

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mpppk/imagine/domain/model"
	"github.com/mpppk/imagine/domain/repository"
)

type Asset struct {
	assetRepository repository.Asset
}

func NewAsset(assetRepository repository.Asset) *Asset {
	return &Asset{
		assetRepository: assetRepository,
	}
}

func (a *Asset) AddImage(ws model.WSName, filePath string) error {
	// FIXME
	if err := a.assetRepository.Init(ws); err != nil {
		return err
	}
	return a.assetRepository.Add(ws, newAssetFromFilePath(filePath))
}

func (a *Asset) ListAsync(ws model.WSName) (<-chan *model.Asset, error) {
	// FIXME
	if err := a.assetRepository.Init(ws); err != nil {
		return nil, err
	}
	return a.assetRepository.ListByAsync(ws, nil, 10) // FIXME
}

// AssignBoundingBox assign bounding box to asset
func (a *Asset) AssignBoundingBox(ws model.WSName, assetId model.AssetID, box *model.BoundingBox) (*model.Asset, error) {
	// FIXME
	if err := a.assetRepository.Init(ws); err != nil {
		return nil, err
	}
	asset, err := a.assetRepository.Get(ws, assetId) // FIXME
	if err != nil {
		return nil, fmt.Errorf("failed to get asset. id: %v: %w", assetId, err)
	}

	var maxId model.BoundingBoxID = 0
	for _, boundingBox := range asset.BoundingBoxes {
		if boundingBox.ID > maxId {
			maxId = boundingBox.ID
		}
	}
	box.ID = maxId + 1
	asset.BoundingBoxes = append(asset.BoundingBoxes, box)
	if err := a.assetRepository.Update(ws, asset); err != nil {
		return nil, fmt.Errorf("failed to update asset. asset: %#v: %w", asset, err)
	}
	return asset, nil
}

func (a *Asset) UnAssignBoundingBox(ws model.WSName, assetId model.AssetID, boxID model.BoundingBoxID) (*model.Asset, error) {
	// FIXME
	if err := a.assetRepository.Init(ws); err != nil {
		return nil, err
	}
	asset, err := a.assetRepository.Get(ws, assetId) // FIXME
	if err != nil {
		return nil, fmt.Errorf("failed to get asset. id: %v: %w", assetId, err)
	}

	var newBoxes []*model.BoundingBox
	for _, boundingBox := range asset.BoundingBoxes {
		if boundingBox.ID == boxID {
			continue
		}
		newBoxes = append(newBoxes, boundingBox)
	}

	asset.BoundingBoxes = newBoxes
	if err := a.assetRepository.Update(ws, asset); err != nil {
		return nil, fmt.Errorf("failed to update asset. asset: %#v: %w", asset, err)
	}
	return asset, nil
}

func newAssetFromFilePath(filePath string) *model.Asset {
	name := strings.Replace(filepath.Base(filePath), filepath.Ext(filePath), "", -1)
	return &model.Asset{
		Name: name,
		Path: filePath,
	}
}
