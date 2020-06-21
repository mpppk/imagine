//+build wireinject

package registry

//go:generate wire

import (
	bolt "go.etcd.io/bbolt"

	"github.com/google/wire"
	"github.com/labstack/echo"
	"github.com/mpppk/imagine/domain/model"
	"github.com/mpppk/imagine/handler"
	"github.com/mpppk/imagine/infra"
	"github.com/mpppk/imagine/infra/repoimpl"
	"github.com/mpppk/imagine/usecase"
)

// InitializeSumUseCase initialize sum use case with  memorySumHistoryRepository
func InitializeSumUseCase(v []*model.SumHistory) *usecase.Sum {
	wire.Build(repoimpl.NewMemorySumHistory, usecase.NewSum)
	return &usecase.Sum{}
}

// InitializeServer initialize echo server with memorySumHistoryRepository
func InitializeServer(v []*model.SumHistory, b *bolt.DB) *echo.Echo {
	wire.Build(
		handler.New,
		repoimpl.NewMemorySumHistory,
		repoimpl.NewBBoltAsset,
		usecase.NewSum,
		usecase.NewAsset,
		infra.NewServer,
	)
	return &echo.Echo{}
}
