package usecasetest

import (
	"testing"

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

func (t *TagRepository) Add(ws model.WSName, tagWithIndex *model.TagWithIndex) model.TagID {
	tag, err := t.repository.Add(ws, tagWithIndex)
	if err != nil {
		t.t.Fatalf("failed to add tag: %v: %v", err, tagWithIndex)
	}
	return tag
}

type Client struct {
	Asset *AssetRepository
	Tag   *TagRepository
}

func newClient(t *testing.T, c *repository.Client) *Client {
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

func (t *Tag) PutTags(ws model.WSName, tags []*model.Tag) {
	t.t.Helper()
	if err := t.tag.PutTags(ws, tags); err != nil {
		t.t.Fatalf("failed to set tags: %v, %v", err, tags)
	}
}

// SetTags is wrapper for interactor.Tag.SetTags.
func (t *Tag) SetTags(ws model.WSName, tagNames []string) []model.TagID {
	t.t.Helper()
	idList, err := t.tag.SetTags(ws, tagNames)
	if err != nil {
		t.t.Fatalf("failed to set tags: %v, %v", err, tagNames)
	}
	return idList
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
