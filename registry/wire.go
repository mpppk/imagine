//+build wireinject

package registry

import (
	"github.com/google/wire"
	"github.com/mpppk/imagine/action"
	"github.com/mpppk/imagine/domain/repository"
	"github.com/mpppk/imagine/infra"
	"github.com/mpppk/imagine/infra/repoimpl"
	"github.com/mpppk/imagine/usecase/interactor"
	"go.etcd.io/bbolt"
)

//go:generate wire

func NewBoltHandlerCreator(b *bbolt.DB) *action.HandlerCreator {
	wire.Build(
		action.NewHandlerCreator,

		interactor.NewAsset,
		repoimpl.NewBBoltAsset,

		interactor.NewTag,
		repoimpl.NewBBoltTag,

		repoimpl.NewBBoltWorkSpace,
		repoimpl.NewBoltMeta,
		repository.NewClient,
	)
	return nil
}

func InitializeAssetUseCase(b *bbolt.DB) *interactor.Asset {
	wire.Build(
		interactor.NewAsset,
		repoimpl.NewBBoltAsset,
		repoimpl.NewBBoltTag,
	)
	return nil
}

func InitializeTagUseCase(b *bbolt.DB) *interactor.Tag {
	wire.Build(
		interactor.NewTag,
		repoimpl.NewBBoltTag,
	)
	return nil
}

func NewBoltClient(b *bbolt.DB) *repository.Client {
	wire.Build(
		repository.NewClient,
		repoimpl.NewBBoltAsset,
		repoimpl.NewBBoltTag,
		repoimpl.NewBBoltWorkSpace,
		repoimpl.NewBoltMeta,
	)
	return nil
}

func NewBoltUseCases(b *bbolt.DB) *interactor.UseCases {
	wire.Build(
		interactor.New,
		repoimpl.NewBBoltAsset,
		repoimpl.NewBBoltTag,
		repoimpl.NewBBoltWorkSpace,
		repoimpl.NewBoltMeta,
	)
	return nil
}

func NewBoltUseCasesWithDBPath(dbPath string) (*interactor.UseCases, error) {
	wire.Build(
		infra.NewBoltDB,
		interactor.New,
		repoimpl.NewBBoltAsset,
		repoimpl.NewBBoltTag,
		repoimpl.NewBBoltWorkSpace,
		repoimpl.NewBoltMeta,
	)
	return nil, nil
}
