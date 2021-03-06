package interactor

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/mpppk/imagine/usecase"

	"github.com/mpppk/imagine/domain/client"
	"github.com/mpppk/imagine/domain/query"
	"github.com/mpppk/imagine/domain/service/assetsvc"

	"github.com/briandowns/spinner"
	"github.com/mpppk/imagine/domain/model"
	"github.com/mpppk/imagine/domain/repository"
)

type Asset struct {
	assetRepository repository.Asset
	tagRepository   repository.Tag
	tagQuery        query.Tag
}

func NewAsset(assetRepository repository.Asset, tag *client.Tag) usecase.Asset {
	return &Asset{
		assetRepository: assetRepository,
		tagRepository:   tag.TagRepository,
		tagQuery:        tag.TagQuery,
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

func (a *Asset) ReadImportAssetsWithProgressBar(ws model.WSName, reader io.Reader, capacity int, f func(assets []*model.ImportAsset) error) error {
	scanner := bufio.NewScanner(reader)
	cnt := 0
	s := spinner.New(spinner.CharSets[43], 100*time.Millisecond)
	s.Prefix = "loading... "
	s.Start()

	importAssets := make([]*model.ImportAsset, 0, capacity)
	for scanner.Scan() {
		var asset model.ImportAsset
		if err := json.Unmarshal(scanner.Bytes(), &asset); err != nil {
			return fmt.Errorf("failed to unmarshal json to asset: %w", err)
		}
		importAssets = append(importAssets, &asset)
		cnt++
		s.Suffix = strconv.Itoa(cnt)

		if len(importAssets) >= capacity {
			s.Suffix += "(writing...)"
			if err := f(importAssets); err != nil {
				return err
			}
			importAssets = make([]*model.ImportAsset, 0, capacity)
		}
	}

	if len(importAssets) > 0 {
		s.Suffix += "(writing...)"
		if err := f(importAssets); err != nil {
			return err
		}
	}
	s.Stop()
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("faield to scan asset op: %w", err)
	}
	return nil
}

// SaveImportAssetsFromReader updates by Id or path that provided from reader.
// If ID is specified, find asset by ID and update properties. (includes path if it specified)
// If ID is not specified and path is specified, find asset by path and update properties.
// Specified properties are updated and omitted properties are reserved.
// The number of assets specified by capacity is buffered and processed at one time.
// queries can be nil. Only assets that match queries will be merged.
func (a *Asset) SaveImportAssetsFromReader(ws model.WSName, reader io.Reader, capacity int, queries []*model.Query) error {
	f := func(importAssets []*model.ImportAsset) error {
		log.Printf("debug: new batch: %d assets are loaded from reader", capacity)
		if _, _, _, err := a.SaveImportAssets(ws, importAssets, queries); err != nil {
			return fmt.Errorf("faled to add or update assets from reader: %w", err)
		}
		return nil
	}
	return a.ReadImportAssetsWithProgressBar(ws, reader, capacity, f)
}

// SaveImportAssets updates assets by ID or path.
// If ID is specified, find asset by ID and update properties. (includes path if it specified)
// If ID is not specified and path is specified, find asset by path and update properties.
// Specified properties are updated and omitted properties are reserved.
// queries can be nil. Only assets that match queries will be merged.
func (a *Asset) SaveImportAssets(ws model.WSName, importAssets []*model.ImportAsset, queries []*model.Query) (addedAssets, updatedAssets, skippedAssets []*model.Asset, err error) {
	errMsg := "failed to save import assets"
	if _, err := a.tagRepository.AddByNames(ws, assetsvc.ToUniqTagNames(importAssets)); err != nil {
		return nil, nil, nil, fmt.Errorf("%s: %w", errMsg, err)
	}

	tagSet, err := a.tagQuery.ListAsSet(ws)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s: %w", errMsg, err)
	}

	assets, err := assetsvc.ToAssets(importAssets, tagSet)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s: %w", errMsg, err)
	}

	assetsWithID, assetsWithOutID := assetsvc.SplitByID(assets)
	assetsWithPath, assetsWithOutIDAndPath := assetsvc.SplitByPath(assetsWithOutID)

	needToAddAssets := assetsWithOutIDAndPath

	updatedAssets1, filteredAssets1, skippedAssets1, err := a.BatchUpdateByID(ws, assetsWithID, queries)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s: %w", errMsg, err)
	}
	needToAddAssets = append(needToAddAssets, skippedAssets...)

	updatedAssets2, filteredAssets2, skippedAssets2, err := a.BatchUpdateByPath(ws, assetsWithPath, queries)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s: %w", errMsg, err)
	}

	needToAddAssets = append(needToAddAssets, skippedAssets1...)
	needToAddAssets = append(needToAddAssets, skippedAssets2...)

	addedIDList, err := a.assetRepository.BatchAdd(ws, needToAddAssets)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s: %w", errMsg, err)
	}

	if err := assetsvc.ReRegister(needToAddAssets, addedIDList); err != nil {
		return nil, nil, nil, fmt.Errorf("%s: %w", errMsg, err)
	}
	return needToAddAssets, append(updatedAssets1, updatedAssets2...), append(filteredAssets1, filteredAssets2...), nil
}

func (a *Asset) queryIndex(ws model.WSName, assets []*model.Asset, queries []*model.Query, matchToNil bool) (matchedAssetsIndex, filteredAssetsIndex []int, err error) {
	errMsg := "failed to query assets"
	tagSet, err := a.tagQuery.ListAsSet(ws)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", errMsg, err)
	}
	matchedAssetsIndex, filteredAssetsIndex = assetsvc.QueryIndex(assets, queries, tagSet, matchToNil)
	return matchedAssetsIndex, filteredAssetsIndex, nil
}

// BatchUpdateByID update by provided assets.
// If DB already has same ID asset, merge it and provided asset, then save the asset.
// If DB does not have same ID asset, do nothing and return as skippedAssets.
// queries can be nil. Only assets that match queries will be merged. Unmatched assets will be returned as filteredAssets.
func (a *Asset) BatchUpdateByID(ws model.WSName, assets []*model.Asset, queries []*model.Query) (updatedAssets, filteredAssets, skippedAssets []*model.Asset, err error) {
	errMsg := "failed to batch merge by ID"

	assetIDList := assetsvc.ToAssetIDList(assets)
	newAssets, err := a.assetRepository.ListByIDList(ws, assetIDList)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s: %w", errMsg, err)
	}

	skippedAssets = assetsvc.FilterByIndex(assets, assetsvc.NilIndex(newAssets))

	matchedAssetsIndex, filteredAssetsIndex, err := a.queryIndex(ws, newAssets, queries, true)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s: %w", errMsg, err)
	}

	if err := assetsvc.Update(newAssets, assets); err != nil {
		return nil, nil, nil, fmt.Errorf("%s: %w", errMsg, err)
	}

	matchedNewAssets := assetsvc.FilterByIndex(newAssets, matchedAssetsIndex)
	filteredNewAssets := assetsvc.FilterByIndex(assets, filteredAssetsIndex)

	updatedAssets, _, err = a.assetRepository.BatchUpdateByID(ws, matchedNewAssets)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to update assets by ID: %w", err)
	}

	return updatedAssets, filteredNewAssets, skippedAssets, nil
}

// BatchUpdateByPath update by provided assets.
// If DB already has same ID asset, merge it and provided asset, then save the asset.
// If DB does not have same ID asset, do nothing and return as skippedAssets.
// queries can be nil. Only assets that match queries will be merged.
// Note: filteredAssets will be returned with updated properties, not original.
func (a *Asset) BatchUpdateByPath(ws model.WSName, assets []*model.Asset, queries []*model.Query) (updatedAssets, filteredAssets, skippedAssets []*model.Asset, err error) {
	errMsg := "failed to batch merge by path"

	assetPaths := assetsvc.ToPaths(assets)
	newAssets, err := a.assetRepository.ListByPaths(ws, assetPaths)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s: %w", errMsg, err)
	}

	for i, newAsset := range newAssets {
		if newAsset == nil {
			skippedAssets = append(skippedAssets, assets[i])
		}
	}

	matchedAssetsIndex, filteredAssetsIndex, err := a.queryIndex(ws, newAssets, queries, true)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s: %w", errMsg, err)
	}

	if err := assetsvc.Update(newAssets, assets); err != nil {
		return nil, nil, nil, fmt.Errorf("%s: %w", errMsg, err)
	}
	matchedNewAssets := assetsvc.FilterByIndex(newAssets, matchedAssetsIndex)
	matchedNewAssets = assetsvc.FilterNil(matchedNewAssets)
	filteredNewAssets := assetsvc.FilterByIndex(assets, filteredAssetsIndex)

	updatedAssets, skippedAssets2, err := a.assetRepository.BatchUpdateByPath(ws, matchedNewAssets)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to update assets by ID: %w", err)
	}
	skippedAssets = append(skippedAssets, skippedAssets2...)
	return updatedAssets, filteredNewAssets, skippedAssets, nil
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

func (a *Asset) ListAsyncByQueries(ctx context.Context, ws model.WSName, queries []*model.Query) (<-chan *model.Asset, error) {
	tagSet, err := a.tagQuery.ListAsSet(ws)
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

func (a *Asset) ListAsyncWithFormat(wsName model.WSName, formatType string, capacity int) (<-chan string, <-chan error, error) {
	assetChan, err := a.assetRepository.ListByAsync(context.Background(), wsName, nil, capacity)
	if err != nil {
		return nil, nil, err
	}

	tagSet, err := a.tagQuery.ListAsSet(wsName)
	if err != nil {
		return nil, nil, err
	}

	format := func(format string, asset *model.Asset) (string, error) {
		switch format {
		case "json":
			return asset.ToJson()
		case "csv":
			return asset.ToCSVRow(tagSet)
		default:
			return "", fmt.Errorf("unknown output format: %s", format)
		}
	}

	outCh := make(chan string, capacity)
	errCh := make(chan error, 1)

	go func() {
		if formatType == "csv" {
			header := []string{strconv.Quote("id"), strconv.Quote("path"), strconv.Quote("tags")}
			outCh <- strings.Join(header, ",")
		}

		for asset := range assetChan {
			t, err := format(formatType, asset)
			if err != nil {
				errCh <- err
				return
			}
			outCh <- t
		}
		close(outCh)
	}()
	return outCh, errCh, nil
}
