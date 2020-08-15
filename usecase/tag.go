package usecase

import (
	"github.com/mpppk/imagine/domain/model"
	"github.com/mpppk/imagine/domain/repository"
)

type Tag struct {
	tagRepository repository.Tag
}

func NewTag(tagRepository repository.Tag) *Tag {
	return &Tag{
		tagRepository: tagRepository,
	}
}

func (a *Tag) List(ws model.WSName) ([]*model.Tag, error) {
	// FIXME
	if err := a.tagRepository.Init(ws); err != nil {
		return nil, err
	}
	return a.tagRepository.ListAll(ws)
}
