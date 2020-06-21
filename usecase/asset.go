package usecase

import "github.com/mpppk/imagine/domain/repository"

type Asset struct {
	assetRepository repository.Asset
}

func NewAsset(assetRepository repository.Asset) *Asset {
	return &Asset{
		assetRepository: assetRepository,
	}
}

func (a *Asset) AddImages(filePaths []string) {

}
