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
