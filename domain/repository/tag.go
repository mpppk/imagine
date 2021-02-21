//go:generate mockgen -source=./tag.go -destination=../../infra/repoimpl/mock_tag.go -package=repoimpl

package repository

import "github.com/mpppk/imagine/domain/model"

type Tag interface {
	Init(ws model.WSName) error
	Add(ws model.WSName, tagWithIndex *model.TagWithIndex) (model.TagID, error)
	AddByName(ws model.WSName, tagName string) (model.TagID, bool, error)
	AddByNames(ws model.WSName, tagNames []string) ([]model.TagID, error)
	Get(ws model.WSName, id model.TagID) (tagWithIndex *model.TagWithIndex, exist bool, err error)

	// UpdateByTagSet update tags by provided TagSet.
	// If tag exists on TagSet, but is not persisted yet, add the tag to storage.
	// If tag exists on both TagSet and storage, do nothing or update properties.
	// If tag does not exist on TagSet, but exists on storage, delete the tag from storage.
	// In either case, tag identity is determined using the ID.
	// Note: TagRepository implementation must ensure that if you add a previously deleted tag again, the tag will have the same ID as last time.
	//UpdateByTagSet(ws model.WSName, set *model.TagSet) error

	// Save persists tag and return ID of the tag.
	// If the tag already exists, update it. Otherwise add new tag with new ID.
	Save(ws model.WSName, tagWithIndex *model.TagWithIndex) (model.TagID, error)

	RecreateBucket(ws model.WSName) error
	ListAll(ws model.WSName) ([]*model.TagWithIndex, error)
	ListByAsync(ws model.WSName, f func(tagWithIndex *model.TagWithIndex) bool, cap int) (tagChan <-chan *model.TagWithIndex, err error)
	ListBy(ws model.WSName, f func(tagWithIndex *model.TagWithIndex) bool) (tags []*model.TagWithIndex, err error)
	ListAsSet(ws model.WSName) (set *model.TagSet, err error)
	ForEach(ws model.WSName, f func(tagWithIndex *model.TagWithIndex) error) error

	// Delete deletes tag which have given ID
	Delete(ws model.WSName, idList []model.TagID) error
}
