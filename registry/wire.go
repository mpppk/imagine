//+build wireinject

package registry

//go:generate wire

import (
	"github.com/google/wire"
	"github.com/mpppk/imagine/action"
	"github.com/mpppk/imagine/infra/repoimpl"
	"github.com/mpppk/imagine/usecase"
	"go.etcd.io/bbolt"
)

func InitializeDirectoryScanHandler(b *bbolt.DB) *action.DirectoryScanHandler {
	wire.Build(usecase.NewAsset, repoimpl.NewBBoltAsset, action.NewReadDirectoryScanHandler)
	return &action.DirectoryScanHandler{}
}

func InitializeRequestWorkSpacesHandler(b *bbolt.DB) *action.RequestWorkSpacesHandler {
	wire.Build(repoimpl.NewBBoltGlobal, action.NewRequestWorkSpacesHandler)
	return nil
}

func InitializeRequestAssetsHandler(b *bbolt.DB) *action.RequestAssetsHandler {
	wire.Build(usecase.NewAsset, repoimpl.NewBBoltAsset, action.NewRequestAssetsHandler)
	return &action.RequestAssetsHandler{}
}
