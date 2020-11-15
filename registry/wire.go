//+build wireinject

package registry

import (
	"github.com/google/wire"
	"github.com/mpppk/imagine/action"
	"github.com/mpppk/imagine/infra/repoimpl"
	"github.com/mpppk/imagine/usecase"
	"go.etcd.io/bbolt"
)

//go:generate wire

func InitializeHandlerCreator(b *bbolt.DB) *action.HandlerCreator {
	wire.Build(
		action.NewHandlerCreator,

		usecase.NewAsset,
		repoimpl.NewBBoltAsset,

		usecase.NewTag,
		repoimpl.NewBBoltTag,

		repoimpl.NewBBoltWorkSpace,
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
