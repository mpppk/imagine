//go:generate mockgen -source=./tag.go -destination=../../infra/queryimpl/mock_tag.go -package=queryimpl

package query

import "github.com/mpppk/imagine/domain/model"

type Tag interface {
	Init(ws model.WSName) error
	Get(ws model.WSName, id model.TagID) (tagWithIndex *model.TagWithIndex, exist bool, err error)
	ListAll(ws model.WSName) ([]*model.TagWithIndex, error)
	ListByAsync(ws model.WSName, f func(tagWithIndex *model.TagWithIndex) bool, cap int) (tagChan <-chan *model.TagWithIndex, err error)
	ListBy(ws model.WSName, f func(tagWithIndex *model.TagWithIndex) bool) (tags []*model.TagWithIndex, err error)
	ListAsSet(ws model.WSName) (set *model.TagSet, err error)
	ForEach(ws model.WSName, f func(tagWithIndex *model.TagWithIndex) error) error
}
