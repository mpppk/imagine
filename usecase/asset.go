package usecase

import (
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

func newAssetFromFilePath(filePath string) *model.Asset {
	name := strings.Replace(filepath.Base(filePath), filepath.Ext(filePath), "", -1)
	return &model.Asset{
		Name: name,
		Path: filePath,
	}
}
