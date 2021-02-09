//go:generate mockgen -source=./tag.go -destination=./mock_usecase/mock_tag.go

package usecase

import "github.com/mpppk/imagine/domain/model"

type Tag interface {
	List(ws model.WSName) (tags []*model.Tag, err error)
	SetTags(ws model.WSName, tags []*model.Tag) error
}
