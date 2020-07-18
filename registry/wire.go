//+build wireinject

package registry

//go:generate wire

import (
	"github.com/google/wire"
	"github.com/mpppk/imagine/domain/model"
	"github.com/mpppk/imagine/infra/repoimpl"
	"github.com/mpppk/imagine/usecase"
)

// InitializeSumUseCase initialize sum use case with  memorySumHistoryRepository
func InitializeSumUseCase(v []*model.SumHistory) *usecase.Sum {
	wire.Build(repoimpl.NewMemorySumHistory, usecase.NewSum)
	return &usecase.Sum{}
}
