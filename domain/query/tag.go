//go:generate mockgen -source=./tag.go -destination=../../infra/queryimpl/mock_tag.go -package=queryimpl

package query

import "github.com/mpppk/imagine/domain/model"

type Tag interface {
	Init(ws model.WSName) error
	Get(ws model.WSName, id model.TagID) (tagWithIndex *model.Tag, exist bool, err error)
	ListAll(ws model.WSName) ([]*model.Tag, error)
	ListByAsync(ws model.WSName, f func(tagWithIndex *model.Tag) bool, cap int) (tagChan <-chan *model.Tag, err error)
	ListBy(ws model.WSName, f func(tagWithIndex *model.Tag) bool) (tags []*model.Tag, err error)
	ListAsSet(ws model.WSName) (set *model.TagSet, err error)
	ForEach(ws model.WSName, f func(tagWithIndex *model.Tag) error) error
}
