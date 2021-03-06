//+build wireinject

package registry

import (
	"github.com/google/wire"
	"github.com/mpppk/imagine/action"
	"github.com/mpppk/imagine/domain/client"
	"github.com/mpppk/imagine/infra"
	"github.com/mpppk/imagine/infra/queryimpl"
	"github.com/mpppk/imagine/infra/repoimpl"
	"github.com/mpppk/imagine/usecase"
	"github.com/mpppk/imagine/usecase/interactor"
	"go.etcd.io/bbolt"
)

//go:generate wire

func NewBoltHandlerCreator(b *bbolt.DB) *action.HandlerCreator {
	wire.Build(
		action.NewHandlerCreator,

		interactor.NewAsset,
		repoimpl.NewBBoltAsset,

		client.NewTag,
		interactor.NewTag,
		repoimpl.NewBBoltTag,
		queryimpl.NewBBoltTag,

		repoimpl.NewBBoltWorkSpace,
		repoimpl.NewBoltMeta,
		client.New,
	)
	return nil
}

func InitializeAssetUseCase(b *bbolt.DB) usecase.Asset {
	wire.Build(
		client.NewTag,
		interactor.NewAsset,
		repoimpl.NewBBoltAsset,
		queryimpl.NewBBoltTag,
		repoimpl.NewBBoltTag,
	)
	return nil
}

func InitializeTagUseCase(b *bbolt.DB) *interactor.Tag {
	wire.Build(
		client.NewTag,
		interactor.NewTag,
		queryimpl.NewBBoltTag,
		repoimpl.NewBBoltTag,
	)
	return nil
}

func NewBoltClient(b *bbolt.DB) *client.Client {
	wire.Build(
		client.New,
		client.NewTag,
		queryimpl.NewBBoltTag,
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
		client.NewTag,
		queryimpl.NewBBoltTag,
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
		client.NewTag,
		queryimpl.NewBBoltTag,
		repoimpl.NewBBoltAsset,
		repoimpl.NewBBoltTag,
		repoimpl.NewBBoltWorkSpace,
		repoimpl.NewBoltMeta,
	)
	return nil, nil
}
