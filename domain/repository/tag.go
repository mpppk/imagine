package repository

import "github.com/mpppk/imagine/domain/model"

type Tag interface {
	Init(ws model.WSName) error
	Add(ws model.WSName, tag *model.Tag) (model.TagID, error)
	Get(ws model.WSName, id model.TagID) (tag *model.Tag, exist bool, err error)
	Update(ws model.WSName, tag *model.Tag) error
	RecreateBucket(ws model.WSName) error
	ListAll(ws model.WSName) ([]*model.Tag, error)
	ListByAsync(ws model.WSName, f func(tag *model.Tag) bool, cap int) (tagChan <-chan *model.Tag, err error)
	ListBy(ws model.WSName, f func(tag *model.Tag) bool) (tags []*model.Tag, err error)
	ForEach(ws model.WSName, f func(tag *model.Tag) error) error
}
