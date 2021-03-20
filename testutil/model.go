package testutil

import (
	"sort"

	"github.com/mpppk/imagine/domain/model"
)

// MustNewUnindexedTag construct and returns UnindexedTag.
// Panic if invalid parameters are provided.
func MustNewUnindexedTag(id model.TagID, name string) *model.UnindexedTag {
	tag, err := model.NewUnindexedTag(id, name)
	PanicIfErrExist(err)
	return tag
}

// MustNewUnregisteredUnindexedTag construct and returns UnindexedTag.
// Panic if invalid parameters are provided.
func MustNewUnregisteredUnindexedTag(name string) *model.UnregisteredUnindexedTag {
	tag, err := model.NewUnregisteredUnindexedTag(name)
	PanicIfErrExist(err)
	return tag
}

// MustNewTag construct and returns Tag.
// Panic if invalid parameters are provided.
func MustNewTag(id model.TagID, name string, index int) *model.Tag {
	tag, err := model.NewTag(id, name, index)
	PanicIfErrExist(err)
	return tag
}

// MustNewUnregisteredTag construct and returns Tag.
// Panic if invalid parameters are provided.
func MustNewUnregisteredTag(name string, index int) *model.UnregisteredTag {
	tag, err := model.NewUnregisteredTag(name, index)
	PanicIfErrExist(err)
	return tag
}

// ReadAllAssetsFromCh read all assets from channel
func ReadAllAssetsFromCh(ch <-chan *model.Asset) (assets []*model.Asset) {
	for asset := range ch {
		assets = append(assets, asset)
	}
	return
}

func SortBoundingBoxesByTagID(asset *model.Asset) {
	sort.Slice(asset.BoundingBoxes, func(i, j int) bool {
		return asset.BoundingBoxes[i].TagID < asset.BoundingBoxes[j].TagID
	})
}

func SortTagsByID(asset *model.Asset) {
	sort.Slice(asset.BoundingBoxes, func(i, j int) bool {
		return asset.BoundingBoxes[i].TagID < asset.BoundingBoxes[j].TagID
	})
}
