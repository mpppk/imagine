package usecase

import (
	"fmt"

	"github.com/mpppk/imagine/domain/model"
	"github.com/mpppk/imagine/domain/repository"
)

type UseCases struct {
	Asset     *Asset
	Tag       *Tag
	Migration *Migration
	Client    *repository.Client
}

func New(asset repository.Asset, tag repository.Tag, workspace repository.WorkSpace, meta repository.Meta) *UseCases {
	return &UseCases{
		Asset:     NewAsset(asset, tag),
		Tag:       NewTag(tag),
		Migration: NewMigration(asset, meta),
		Client:    repository.NewClient(asset, tag, workspace, meta),
	}
}

func (u *UseCases) InitializeWorkSpace(ws model.WSName) error {
	if err := u.Asset.Init(ws); err != nil {
		return fmt.Errorf("failed to initialize asset usecase: %w", err)
	}
	return nil
}

func (u *UseCases) Close() error {
	return u.Client.Close()
}
