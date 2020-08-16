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
	}
}

func (h *HandlerCreator) NewDirectoryScanHandler() *FSScanHandler {
	return NewFSScanHandler(h.assetUseCase)
}

func (h *HandlerCreator) NewRequestWorkSpacesHandler() *RequestWorkSpacesHandler {
	return NewRequestWorkSpacesHandler(h.globalRepository)
}

func (h *HandlerCreator) NewRequestAssetsHandler() *RequestAssetsHandler {
	return NewRequestAssetsHandler(h.assetUseCase)
}

func (h *HandlerCreator) NewTagRequestHandler() *TagRequestHandler {
	return NewTagRequestHandler(h.tagUseCase)
}

func (h *HandlerCreator) NewTagUpdateHandler() *TagUpdateHandler {
	return NewTagUpdateHandler(h.tagUseCase)
}

func (h *HandlerCreator) NewTagUpdate() *TagUpdateHandler {
	return NewTagUpdateHandler(h.tagUseCase)
}
