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
	if err := a.tagRepository.Init(ws); err != nil {
		return err
	}
	return a.assetRepository.Init(ws)
}

type AssetImportResult struct {
	Asset *model.ImportAsset
	Err   error
}

func (a *Asset) ImportBoundingBoxesFromReader(ws model.WSName, reader io.Reader) error {
	scanner := bufio.NewScanner(reader)
	cnt := 0
	cap1 := 10000 // FIXME
	cap2 := 5000  // FIXME
	s := spinner.New(spinner.CharSets[43], 100*time.Millisecond)
	s.Prefix = "loading... "
	s.Start()

	importAssets := make([]*model.ImportAsset, 0, cap1)
	for scanner.Scan() {
		var asset model.ImportAsset
		if err := json.Unmarshal(scanner.Bytes(), &asset); err != nil {
			return fmt.Errorf("failed to unmarshal json to asset")
		}
		importAssets = append(importAssets, &asset)
		cnt++
		s.Suffix = strconv.Itoa(cnt)

		if len(importAssets) >= cap1 {
			s.Suffix += "(writing...)"
			if _, err := a.AppendBoundingBoxes(ws, importAssets, cap2); err != nil {
				return fmt.Errorf("failed to add import assets: %w", err)
			}
			importAssets = make([]*model.ImportAsset, 0, cap1)
		}
	}

	if len(importAssets) > 0 {
		s.Suffix += "(writing...)"
		if _, err := a.AddImportAssets(ws, importAssets, cap2); err != nil {
			return fmt.Errorf("failed to add import assets: %w", err)
		}
	}
	s.Stop()
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("faield to scan asset op: %w", err)
	}
	return nil
}
func (a *Asset) ImportFromReader(ws model.WSName, reader io.Reader, new bool) error {
	scanner := bufio.NewScanner(reader)
	cnt := 0
	cap1 := 10000 // FIXME
	cap2 := 5000  // FIXME
	s := spinner.New(spinner.CharSets[43], 100*time.Millisecond)
	s.Prefix = "loading... "
	s.Start()

	importAssets := make([]*model.ImportAsset, 0, cap1)
	for scanner.Scan() {
		var asset model.ImportAsset
		if err := json.Unmarshal(scanner.Bytes(), &asset); err != nil {
			return fmt.Errorf("failed to unmarshal json to asset")
		}
		importAssets = append(importAssets, &asset)
		cnt++
		s.Suffix = strconv.Itoa(cnt)

		if len(importAssets) >= cap1 {
			s.Suffix += "(writing...)"
			if _, err := a.AddImportAssets(ws, importAssets, cap2); err != nil {
				return fmt.Errorf("failed to add import assets: %w", err)
			}
			importAssets = make([]*model.ImportAsset, 0, cap1)
		}
	}

	if len(importAssets) > 0 {
		s.Suffix += "(writing...)"
		if _, err := a.AddImportAssets(ws, importAssets, cap2); err != nil {
			return fmt.Errorf("failed to add import assets: %w", err)
		}
	}
	s.Stop()
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("faield to scan asset op: %w", err)
	}
	return nil
}

func (a *Asset) AppendBoundingBoxes(ws model.WSName, assets []*model.ImportAsset, cap int) ([]model.AssetID, error) {
	updateAssets := make([]*model.Asset, 0, cap)
	var idList []model.AssetID

	tagSet, err := a.tagRepository.ListAsSet(ws)
	if err != nil {
		return nil, fmt.Errorf("failed to get tag list: %w", err)
	}

	for _, asset := range assets {
		for _, box := range asset.BoundingBoxes {
			if _, ok := tagSet.GetByName(box.TagName); !ok {
				id, _, err := a.tagRepository.AddByName(ws, box.TagName)
				if err != nil {
					return nil, fmt.Errorf("failed to add tag. name is %s: %w", box.TagName, err)
				}
				tagSet.Set(&model.Tag{
					ID:   id,
					Name: box.TagName,
				})
			}
		}

		ast, err := asset.ToAsset(tagSet)
		if err != nil {
			return nil, fmt.Errorf("failed to convert to asset from import asset: %w", err)
		}

		if asset.Path != "" {
			updateAssets = append(updateAssets, ast)
		} else {
			log.Printf("warning: json line is ignored because image path is empty")
		}

		if len(updateAssets) >= cap {
			idl, err := a.assetRepository.BatchAppendBoundingBoxes(ws, updateAssets)
			if err != nil {
				return nil, fmt.Errorf("failed to append bounding boxes: %w", err)
			}

			idList = append(idList, idl...)
			updateAssets = make([]*model.Asset, 0, cap)
		}
	}

	if len(updateAssets) > 0 {
		idl, err := a.assetRepository.BatchAppendBoundingBoxes(ws, updateAssets)
		if err != nil {
			return nil, fmt.Errorf("failed to append bounding boxes: %w", err)
		}

		idList = append(idList, idl...)
	}

	return idList, nil
}
func (a *Asset) AddImportAssets(ws model.WSName, assets []*model.ImportAsset, cap int) ([]model.AssetID, error) {
	newAssets := make([]*model.Asset, 0, cap)
	updateAssets := make([]*model.Asset, 0, cap)
	var idList []model.AssetID

	tagSet, err := a.tagRepository.ListAsSet(ws)
	if err != nil {
		return nil, fmt.Errorf("failed to get tag list: %w", err)
	}

	for _, asset := range assets {
		for _, box := range asset.BoundingBoxes {
			if _, ok := tagSet.GetByName(box.TagName); !ok {
				id, _, err := a.tagRepository.AddByName(ws, box.TagName)
				if err != nil {
					return nil, fmt.Errorf("failed to add tag. name is %s: %w", box.TagName, err)
				}
				tagSet.Set(&model.Tag{
					ID:   id,
					Name: box.TagName,
				})
			}
		}

		ast, err := asset.ToAsset(tagSet)
		if err != nil {
			return nil, fmt.Errorf("failed to convert to asset from import asset: %w", err)
		}

		if asset.ID == 0 {
			if asset.Path == "" {
				log.Printf("warning: image path is empty")
				return nil, nil
			}
			newAssets = append(newAssets, ast)
		} else {
			updateAssets = append(updateAssets, ast)
		}

		if len(newAssets) >= cap {
			idl, err := a.assetRepository.BatchAdd(ws, newAssets)
			if err != nil {
				return nil, fmt.Errorf("failed to add asset from image path: %w", err)
			}
			idList = append(idList, idl...)
			newAssets = make([]*model.Asset, 0, cap)
		}

		if len(updateAssets) >= cap {
			if err := a.assetRepository.BatchUpdate(ws, updateAssets); err != nil {
				return nil, fmt.Errorf("failed to update assets: %w", err)
			}

			for _, updateAsset := range updateAssets {
				idList = append(idList, updateAsset.ID)
			}

			updateAssets = make([]*model.Asset, 0, cap)
		}
	}

	if len(newAssets) > 0 {
		idl, err := a.assetRepository.BatchAdd(ws, newAssets)
		if err != nil {
			return nil, fmt.Errorf("failed to add asset from image path: %w", err)
		}
		idList = append(idList, idl...)
	}

	if len(updateAssets) > 0 {
		if err := a.assetRepository.BatchUpdate(ws, updateAssets); err != nil {
			return nil, fmt.Errorf("failed to update assets: %w", err)
		}

		for _, updateAsset := range updateAssets {
			idList = append(idList, updateAsset.ID)
		}
	}

	return idList, nil
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
