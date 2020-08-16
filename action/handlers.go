package action

import (
	"github.com/mpppk/imagine/domain/repository"
	"github.com/mpppk/imagine/usecase"
	"go.etcd.io/bbolt"
)

type HandlerCreator struct {
	assetUseCase     *usecase.Asset
	tagUseCase       *usecase.Tag
	globalRepository repository.Global
	b                *bbolt.DB
	Asset            *assetHandlerCreator
	Tag              *tagHandlerCreator
	FS               *fsHandlerCreator
}

func NewHandlerCreator(
	assetUseCase *usecase.Asset,
	tagUseCase *usecase.Tag,
	globalRepository repository.Global,
	b *bbolt.DB,
) *HandlerCreator {
	return &HandlerCreator{
		assetUseCase:     assetUseCase,
		tagUseCase:       tagUseCase,
		globalRepository: globalRepository,
		b:                b,
		Asset:            newAssetHandlerCreator(assetUseCase),
		Tag:              newTagHandlerCreator(tagUseCase),
		FS:               newFSHandlerCreator(assetUseCase),
	}
}

func (h *HandlerCreator) NewRequestWorkSpacesHandler() *requestWorkSpacesHandler {
	return &requestWorkSpacesHandler{globalRepository: h.globalRepository}
}
