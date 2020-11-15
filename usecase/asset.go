package usecase

import (
	"context"
	"fmt"

	"github.com/mpppk/imagine/domain/model"
	"github.com/mpppk/imagine/domain/repository"
)

type Asset struct {
	assetRepository repository.Asset
	tagRepository   repository.Tag
}

func NewAsset(assetRepository repository.Asset, tagRepository repository.Tag) *Asset {
	return &Asset{
		assetRepository: assetRepository,
		tagRepository:   tagRepository,
	}
}

func (a *Asset) Init(ws model.WSName) error {
	return a.assetRepository.Init(ws)
}

func (a *Asset) AddAssetFromImagePathListIfDoesNotExist(ws model.WSName, filePathList []string) ([]model.AssetID, error) {
	return a.assetRepository.AddByFilePathListIfDoesNotExist(ws, filePathList)
}

func (a *Asset) AddAssetFromImagePathIfDoesNotExist(ws model.WSName, filePath string) (model.AssetID, bool, error) {
	return a.assetRepository.AddByFilePathIfDoesNotExist(ws, filePath)
}

func (a *Asset) AddAssetFromImagePath(ws model.WSName, filePath string) (model.AssetID, error) {
	// FIXME
	if err := a.assetRepository.Init(ws); err != nil {
		return 0, err
	}
	return a.assetRepository.Add(ws, model.NewAssetFromFilePath(filePath))
}

func (a *Asset) ListAsync(ctx context.Context, ws model.WSName) (<-chan *model.Asset, error) {
	return a.assetRepository.ListByAsync(ctx, ws, nil, 50) // FIXME
}

func (a *Asset) ListAsyncByQueries(ctx context.Context, ws model.WSName, queries []*model.Query) (<-chan *model.Asset, error) {
	tagSet, err := a.tagRepository.ListAsSet(ws)
	if err != nil {
		return nil, err
	}

	if query, ok := hasPathQuery(queries); ok {
		return a.handlePathQuery(ws, queries, query.Value, tagSet)
	}

	f := func(asset *model.Asset) bool {
		return checkQueries(asset, queries, tagSet)
	}
	return a.assetRepository.ListByAsync(ctx, ws, f, 50) // FIXME
}

func (a *Asset) handlePathQuery(ws model.WSName, queries []*model.Query, path string, tagSet *model.TagSet) (<-chan *model.Asset, error) {
	asset, exist, err := a.assetRepository.GetByPath(ws, path)
	if err != nil {
		return nil, fmt.Errorf("failed to handle path query: %w", err)
	}
	c := make(chan *model.Asset, 1)

	if exist && checkQueries(asset, queries, tagSet) {
		c <- asset
	}
	close(c)
	return c, nil
}

func hasPathQuery(queries []*model.Query) (*model.Query, bool) {
	for _, query := range queries {
		if query.Op == model.PathEqualsQueryOP {
			return query, true
		}
	}
	return nil, false
}

func checkQueries(asset *model.Asset, queries []*model.Query, tagSet *model.TagSet) bool {
	for _, query := range queries {
		if !query.Match(asset, tagSet) {
			return false
		}
	}
	return true
}

// AssignBoundingBox assign bounding box to asset
func (a *Asset) AssignBoundingBox(ws model.WSName, assetID model.AssetID, box *model.BoundingBox) (*model.Asset, error) {
	// FIXME
	if err := a.assetRepository.Init(ws); err != nil {
		return nil, err
	}
	asset, _, err := a.assetRepository.Get(ws, assetID) // FIXME
	if err != nil {
		return nil, fmt.Errorf("failed to get asset. id: %v: %w", assetID, err)
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

func (a *Asset) UnAssignBoundingBox(ws model.WSName, assetID model.AssetID, boxID model.BoundingBoxID) (*model.Asset, error) {
	// FIXME
	if err := a.assetRepository.Init(ws); err != nil {
		return nil, err
	}
	asset, _, err := a.assetRepository.Get(ws, assetID) // FIXME
	if err != nil {
		return nil, fmt.Errorf("failed to get asset. id: %v: %w", assetID, err)
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

func (a *Asset) ModifyBoundingBox(ws model.WSName, assetID model.AssetID, box *model.BoundingBox) (*model.Asset, error) {
	// FIXME
	if err := a.assetRepository.Init(ws); err != nil {
		return nil, err
	}
	asset, exist, err := a.assetRepository.Get(ws, assetID) // FIXME
	if err != nil {
		return nil, fmt.Errorf("failed to get asset. id: %v: %w", assetID, err)
	}

	if !exist {
		return nil, fmt.Errorf("asset does not found when modify bounding box. asset id:%v", assetID)
	}

	asset.BoundingBoxes = model.ReplaceBoundingBoxByID(asset.BoundingBoxes, box)
	if err := a.assetRepository.Update(ws, asset); err != nil {
		return nil, fmt.Errorf("failed to update asset. asset: %#v: %w", asset, err)
	}
	return asset, nil
}

func (a *Asset) DeleteBoundingBox(ws model.WSName, assetID model.AssetID, boxID model.BoundingBoxID) error {
	// FIXME
	if err := a.assetRepository.Init(ws); err != nil {
		return err
	}
	asset, _, err := a.assetRepository.Get(ws, assetID) // FIXME
	if err != nil {
		return fmt.Errorf("failed to get asset. id: %v: %w", assetID, err)
	}

	asset.BoundingBoxes = model.RemoveBoundingBoxByID(asset.BoundingBoxes, boxID)
	if err := a.assetRepository.Update(ws, asset); err != nil {
		return fmt.Errorf("failed to update asset. asset: %#v: %w", asset, err)
	}
	return nil
}
