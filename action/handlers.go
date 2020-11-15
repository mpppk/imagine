package action

import (
	"github.com/mpppk/imagine/domain/repository"
	"github.com/mpppk/imagine/usecase"
	"go.etcd.io/bbolt"
)

type HandlerCreator struct {
	assetUseCase *usecase.Asset
	tagUseCase   *usecase.Tag
	b            *bbolt.DB
	Asset        *assetHandlerCreator
	Box          *boxHandlerCreator
	Tag          *tagHandlerCreator
	FS           *fsHandlerCreator
	Workspace    *workspaceHandlerCreator
}

func NewHandlerCreator(
	assetUseCase *usecase.Asset,
	tagUseCase *usecase.Tag,
	client *repository.Client,
	b *bbolt.DB,
) *HandlerCreator {
	return &HandlerCreator{
		assetUseCase: assetUseCase,
		tagUseCase:   tagUseCase,
		b:            b,
		Asset:        newAssetHandlerCreator(assetUseCase),
		Box:          newBoxHandlerCreator(assetUseCase),
		Tag:          newTagHandlerCreator(tagUseCase),
		FS:           newFSHandlerCreator(assetUseCase),
		Workspace:    newWorkspaceHandlerCreator(client),
	}
}
