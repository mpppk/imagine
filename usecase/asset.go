package usecase

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/briandowns/spinner"
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

type AssetImportResult struct {
	Asset *model.ImportAsset
	Err   error
}

func (a *Asset) ImportFromReader(ws model.WSName, reader io.Reader, new bool) error {
	scanner := bufio.NewScanner(reader)
	cnt := 0
	s := spinner.New(spinner.CharSets[43], 100*time.Millisecond)
	s.Prefix = "loading... "
	s.Start()
	for scanner.Scan() {
		var asset model.ImportAsset
		if err := json.Unmarshal(scanner.Bytes(), &asset); err != nil {
			return fmt.Errorf("failed to unmarshal json to asset")
		}
		if _, _, err := a.AddImportAsset(ws, &asset, new); err != nil {
			return fmt.Errorf("failed to import asset: %w", err)
		}
		cnt++
		s.Suffix = strconv.Itoa(cnt)
	}
	s.Stop()
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("faield to scan asset op: %w", err)
	}
	return nil
}

func (a *Asset) AddImportAsset(ws model.WSName, asset *model.ImportAsset, new bool) (model.AssetID, bool, error) {
	if asset.ID == 0 {
		if asset.Path == "" {
			log.Printf("warning: image path is empty")
			return 0, false, nil
		}
		id, err := a.AddAssetFromImagePath(ws, asset.Path)
		if err != nil {
			e := fmt.Errorf("failed to add asset. image path: %s: %w", asset.Path, err)
			return 0, false, e
		}
		log.Printf("debug: asset added: %#v", asset)
		return id, true, nil
	}

	ok, err := a.assetRepository.Has(ws, asset.ID)
	if err != nil {
		e := fmt.Errorf("failed to check asset. image path: %s: %w", asset.Path, err)
		return 0, false, e
	}

	if !ok {
		if new {
			id, err := a.assetRepository.Add(ws, asset.ToAsset())
			if err != nil {
				e := fmt.Errorf("failed to add asset: %w", err)
				return 0, false, e
			}
			log.Printf("debug: asset added: %#v", asset)
			return id, true, nil
		} else {
			log.Printf("debug: asset skipped because it does not exist: id:%d", asset.ID)
			return 0, false, nil
		}
	}

	if err := a.assetRepository.Update(ws, asset.ToAsset()); err != nil {
		e := fmt.Errorf("failed to update asset: %w", err)
		return 0, false, e
	}
	log.Printf("debug: asset updated: %#v", asset)
	return asset.ID, true, nil
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
