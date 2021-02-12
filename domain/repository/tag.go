//go:generate mockgen -source=./tag.go -destination=../../infra/repoimpl/mock_tag.go -package=repoimpl

package repository

import "github.com/mpppk/imagine/domain/model"

type Tag interface {
	Init(ws model.WSName) error
	Add(ws model.WSName, tagWithIndex *model.TagWithIndex) (model.TagID, error)
	AddByName(ws model.WSName, tagName string) (model.TagID, bool, error)
	AddByNames(ws model.WSName, tagNames []string) ([]model.TagID, error)
	Get(ws model.WSName, id model.TagID) (tagWithIndex *model.TagWithIndex, exist bool, err error)

	// Put persists tag.
	// If the tag already exists, update it. Otherwise add new tag.
	Put(ws model.WSName, tagWithIndex *model.TagWithIndex) error

	RecreateBucket(ws model.WSName) error
	ListAll(ws model.WSName) ([]*model.TagWithIndex, error)
	ListByAsync(ws model.WSName, f func(tagWithIndex *model.TagWithIndex) bool, cap int) (tagChan <-chan *model.TagWithIndex, err error)
	ListBy(ws model.WSName, f func(tagWithIndex *model.TagWithIndex) bool) (tags []*model.TagWithIndex, err error)
	ListAsSet(ws model.WSName) (set *model.TagSet, err error)
	ForEach(ws model.WSName, f func(tagWithIndex *model.TagWithIndex) error) error
}
