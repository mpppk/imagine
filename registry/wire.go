//+build wireinject

package registry

import (
	"github.com/google/wire"
	"github.com/mpppk/imagine/action"
	"github.com/mpppk/imagine/domain/repository"
	"github.com/mpppk/imagine/infra"
	"github.com/mpppk/imagine/infra/repoimpl"
	"github.com/mpppk/imagine/usecase"
	"go.etcd.io/bbolt"
)

//go:generate wire

func NewBoltHandlerCreator(b *bbolt.DB) *action.HandlerCreator {
	wire.Build(
		action.NewHandlerCreator,

		usecase.NewAsset,
		repoimpl.NewBBoltAsset,

		usecase.NewTag,
		repoimpl.NewBBoltTag,

		repoimpl.NewBBoltWorkSpace,
		repoimpl.NewBoltMeta,
		repository.NewClient,
	)
	return nil
}

func InitializeAssetUseCase(b *bbolt.DB) *usecase.Asset {
	wire.Build(
		usecase.NewAsset,
		repoimpl.NewBBoltAsset,
		repoimpl.NewBBoltTag,
	)
	return nil
}

func InitializeTagUseCase(b *bbolt.DB) *usecase.Tag {
	wire.Build(
		usecase.NewTag,
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

func NewBoltUseCases(b *bbolt.DB) *usecase.UseCases {
	wire.Build(
		usecase.New,
		repoimpl.NewBBoltAsset,
		repoimpl.NewBBoltTag,
		repoimpl.NewBBoltWorkSpace,
		repoimpl.NewBoltMeta,
	)
	return nil
}

func NewBoltUseCasesWithDBPath(dbPath string) (*usecase.UseCases, error) {
	wire.Build(
		infra.NewBoltDB,
		usecase.New,
		repoimpl.NewBBoltAsset,
		repoimpl.NewBBoltTag,
		repoimpl.NewBBoltWorkSpace,
		repoimpl.NewBoltMeta,
	)
	return nil, nil
}
