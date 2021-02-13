//go:generate mockgen -source=./tag.go -destination=./mock_usecase/mock_tag.go

package usecase

import "github.com/mpppk/imagine/domain/model"

type Tag interface {
	List(ws model.WSName) (tags []*model.Tag, err error)

	// SaveTags persists provided tags.
	// For each tags, if it already exists, update it. Otherwise, add it.
	SaveTags(ws model.WSName, tags []*model.Tag) ([]model.TagID, error)

	// SetTags remove all existing tags and persists provided tags.
	SetTags(ws model.WSName, tagNames []string) ([]model.TagID, error)
}
