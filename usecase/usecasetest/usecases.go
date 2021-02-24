package usecasetest

import (
	"testing"

	"github.com/mpppk/imagine/domain/client"

	"github.com/mpppk/imagine/usecase/interactor"

	"github.com/mpppk/imagine/domain/model"
	"github.com/mpppk/imagine/domain/repository"
)

type AssetRepository struct {
	t          *testing.T
	repository repository.Asset
}

func newAssetRepository(t *testing.T, r repository.Asset) *AssetRepository {
	return &AssetRepository{
		t:          t,
		repository: r,
	}
}

func (a *AssetRepository) List(ws model.WSName) (assets []*model.Asset) {
	a.t.Helper()
	assets, err := a.repository.List(ws)
	if err != nil {
		a.t.Fatalf("failed to list assets: %v", err)
	}
	return assets
}

func (a *AssetRepository) BatchAdd(ws model.WSName, assets []*model.Asset) []model.AssetID {
	a.t.Helper()
	idList, err := a.repository.BatchAdd(ws, assets)
	if err != nil {
		a.t.Fatalf("failed to add assets: %v", err)
	}
	return idList
}

func (a *AssetRepository) ListBy(ws model.WSName, f func(asset *model.Asset) bool) (assets []*model.Asset) {
	a.t.Helper()
	assets, err := a.repository.ListBy(ws, f)
	if err != nil {
		a.t.Fatalf("failed to list assets: %v", err)
	}
	return assets
}

type TagRepository struct {
	t          *testing.T
	repository repository.Tag
}

func newTagRepository(t *testing.T, r repository.Tag) *TagRepository {
	return &TagRepository{
		t:          t,
		repository: r,
	}
}

func (t *TagRepository) AddWithIndex(ws model.WSName, tagWithIndex *model.UnregisteredTag) *model.Tag {
	tag, err := t.repository.AddWithIndex(ws, tagWithIndex)
	if err != nil {
		t.t.Fatalf("failed to add tag: %v: %v", err, tagWithIndex)
	}
	return tag
}

type Client struct {
	Asset *AssetRepository
	Tag   *TagRepository
}

func newClient(t *testing.T, c *client.Client) *Client {
	return &Client{
		Asset: newAssetRepository(t, c.Asset),
		Tag:   newTagRepository(t, c.Tag),
	}
}

type Asset struct {
	t     *testing.T
	asset *interactor.Asset
}

func newAsset(t *testing.T, asset *interactor.Asset) *Asset {
	return &Asset{
		t:     t,
		asset: asset,
	}
}

func (a *Asset) AddOrMergeImportAssets(ws model.WSName, assets []*model.ImportAsset) {
	a.t.Helper()
	err := a.asset.AddOrMergeImportAssets(ws, assets)
	if err != nil {
		a.t.Fatalf("failed to add import assets: %v: %v", err, assets)
	}
}

type Tag struct {
	t   *testing.T
	tag *interactor.Tag
}

func newTag(t *testing.T, tag *interactor.Tag) *Tag {
	return &Tag{
		t:   t,
		tag: tag,
	}
}

// SetTagByNames is wrapper for usecase/SetTagByNames.
func (t *Tag) SetTagByNames(ws model.WSName, tagNames []string) []*model.Tag {
	t.t.Helper()
	tags, err := t.tag.SetTagByNames(ws, tagNames)
	if err != nil {
		t.t.Fatalf("failed to set tags: %v, %v", err, tagNames)
	}
	return tags
}

func (t *Tag) List(ws model.WSName) (tags []*model.Tag) {
	t.t.Helper()
	tags, err := t.tag.List(ws)
	if err != nil {
		t.t.Fatalf("failed to list tags: %v, %v", err, tags)
	}
	return tags
}

type UseCases struct {
	Usecases *interactor.UseCases
	Asset    *Asset
	Tag      *Tag
	Client   *Client
}

func NewUseCases(t *testing.T, u *interactor.UseCases) *UseCases {
	return &UseCases{
		Usecases: u,
		Asset:    newAsset(t, u.Asset),
		Tag:      newTag(t, u.Tag),
		Client:   newClient(t, u.Client),
	}
}
