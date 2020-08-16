package action

import (
	"github.com/mpppk/imagine/domain/repository"
	"github.com/mpppk/imagine/usecase"
	"go.etcd.io/bbolt"
)

type HandlerCreator struct {
	assetUseCase       *usecase.Asset
	tagUseCase         *usecase.Tag
	globalRepository   repository.Global
	b                  *bbolt.DB
	assetActionCreator *assetActionCreator
}

func NewHandlerCreator(
	assetUseCase *usecase.Asset,
	tagUseCase *usecase.Tag,
	globalRepository repository.Global,
	b *bbolt.DB,
) *HandlerCreator {
	return &HandlerCreator{
		assetUseCase:       assetUseCase,
		tagUseCase:         tagUseCase,
		globalRepository:   globalRepository,
		b:                  b,
		assetActionCreator: &assetActionCreator{},
	}
}

func (h *HandlerCreator) NewFSScanHandler() *fsScanHandler {
	return &fsScanHandler{assetUseCase: h.assetUseCase}
}

func (h *HandlerCreator) NewRequestWorkSpacesHandler() *requestWorkSpacesHandler {
	return &requestWorkSpacesHandler{globalRepository: h.globalRepository}
}

func (h *HandlerCreator) NewRequestAssetsHandler() *requestAssetsHandler {
	return &requestAssetsHandler{
		assetUseCase:       h.assetUseCase,
		assetActionCreator: h.assetActionCreator,
	}
}

func (h *HandlerCreator) NewTagRequestHandler() *tagRequestHandler {
	return &tagRequestHandler{tagUseCase: h.tagUseCase}
}

func (h *HandlerCreator) NewTagUpdateHandler() *tagUpdateHandler {
	return &tagUpdateHandler{tagUseCase: h.tagUseCase}
}
