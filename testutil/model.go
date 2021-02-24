package testutil

import (
	"github.com/mpppk/imagine/domain/model"
)

// MustNewTag construct and returns UnindexedTag.
// Panic if invalid parameters are provided.
func MustNewTag(id model.TagID, name string) *model.UnindexedTag {
	tag, err := model.NewTag(id, name)
	PanicIfErrExist(err)
	return tag
}

// MustNewUnregisteredTag construct and returns UnindexedTag.
// Panic if invalid parameters are provided.
func MustNewUnregisteredTag(name string) *model.UnregisteredTag {
	tag, err := model.NewUnregisteredTag(name)
	PanicIfErrExist(err)
	return tag
}

// MustNewTagWithIndex construct and returns TagWithIndex.
// Panic if invalid parameters are provided.
func MustNewTagWithIndex(id model.TagID, name string, index int) *model.TagWithIndex {
	tag, err := model.NewTagWithIndex(id, name, index)
	PanicIfErrExist(err)
	return tag
}

// MustNewUnregisteredTagWithIndex construct and returns TagWithIndex.
// Panic if invalid parameters are provided.
func MustNewUnregisteredTagWithIndex(name string, index int) *model.UnregisteredTagWithIndex {
	tag, err := model.NewUnregisteredTagWithIndex(name, index)
	PanicIfErrExist(err)
	return tag
}
