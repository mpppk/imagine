//go:generate mockgen -source=./tag.go -destination=../../infra/repoimpl/mock_tag.go -package=repoimpl

package repository

import "github.com/mpppk/imagine/domain/model"

type Tag interface {
	Init(ws model.WSName) error
	Add(ws model.WSName, tag *model.UnregisteredUnindexedTag) (*model.Tag, error)
	AddWithIndex(ws model.WSName, unregisteredTag *model.UnregisteredTag) (*model.Tag, error)
	AddByName(ws model.WSName, tagName string) (*model.Tag, bool, error)
	AddByNames(ws model.WSName, tagNames []string) ([]*model.Tag, error)

	// Save persists tag and return ID of the tag.
	// If the tag already exists, update it. Otherwise add new tag with new ID.
	Save(ws model.WSName, tagWithIndex *model.Tag) (*model.Tag, error)

	// Delete deletes tag which have given ID
	Delete(ws model.WSName, idList []model.TagID) error
}
