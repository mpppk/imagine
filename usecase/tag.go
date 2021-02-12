//go:generate mockgen -source=./tag.go -destination=./mock_usecase/mock_tag.go

package usecase

import "github.com/mpppk/imagine/domain/model"

type Tag interface {
	List(ws model.WSName) (tags []*model.Tag, err error)

	// PutTags persists provided tags.
	// For each tags, if it already exists, update it. Otherwise, add it.
	PutTags(ws model.WSName, tags []*model.Tag) error
}
