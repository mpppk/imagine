package interactor

import (
	"fmt"

	"github.com/mpppk/imagine/domain/client"

	"github.com/mpppk/imagine/domain/model"
	"github.com/mpppk/imagine/domain/repository"
)

type UseCases struct {
	Asset     *Asset
	Tag       *Tag
	Migration *Migration
	Client    *client.Client
}

func New(asset repository.Asset, tag *client.Tag, workspace repository.WorkSpace, meta repository.Meta) *UseCases {
	return &UseCases{
		Asset:     NewAsset(asset, tag),
		Tag:       NewTag(tag),
		Migration: NewMigration(asset, meta),
		Client:    client.New(asset, tag, workspace, meta),
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
