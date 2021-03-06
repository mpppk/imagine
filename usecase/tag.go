//go:generate mockgen -source=./tag.go -destination=./mock_usecase/mock_tag.go

package usecase

import "github.com/mpppk/imagine/domain/model"

type Tag interface {
	List(ws model.WSName) (tags []*model.Tag, err error)

	// SaveTags persists provided tags.
	// For each tags, if it already exists, update it. Otherwise, add it.
	SaveTags(ws model.WSName, tags []*model.Tag) ([]*model.Tag, error)

	SetTags(ws model.WSName, tags []*model.UnindexedTag) ([]*model.Tag, error)

	// SetTagByNames remove all existing tags and persists provided tags.
	SetTagByNames(ws model.WSName, tagNames []string) ([]*model.Tag, error)
}
