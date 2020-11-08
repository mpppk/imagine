package repository

import "github.com/mpppk/imagine/domain/model"

type Tag interface {
	Init(ws model.WSName) error
	Add(ws model.WSName, tagWithIndex *model.TagWithIndex) (model.TagID, error)
	Get(ws model.WSName, id model.TagID) (tagWithIndex *model.TagWithIndex, exist bool, err error)
	Update(ws model.WSName, tagWithIndex *model.TagWithIndex) error
	RecreateBucket(ws model.WSName) error
	ListAll(ws model.WSName) ([]*model.TagWithIndex, error)
	ListByAsync(ws model.WSName, f func(tagWithIndex *model.TagWithIndex) bool, cap int) (tagChan <-chan *model.TagWithIndex, err error)
	ListBy(ws model.WSName, f func(tagWithIndex *model.TagWithIndex) bool) (tags []*model.TagWithIndex, err error)
	ForEach(ws model.WSName, f func(tagWithIndex *model.TagWithIndex) error) error
}
