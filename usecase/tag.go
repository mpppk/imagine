package usecase

import (
	"fmt"
	"sort"

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

func (a *Tag) List(ws model.WSName) (tags []*model.Tag, err error) {
	tagsWithIndex, err := a.tagRepository.ListAll(ws)
	if err != nil {
		return nil, err
	}

	sort.Slice(tagsWithIndex, func(i, j int) bool {
		return tagsWithIndex[i].Index < tagsWithIndex[j].Index
	})

	for _, tagWithIndex := range tagsWithIndex {
		tags = append(tags, tagWithIndex.Tag)
	}

	return
}

// SetTags set tag list to workspace
func (a *Tag) SetTags(ws model.WSName, tags []*model.Tag) error {
	// FIXME
	if err := a.tagRepository.Init(ws); err != nil {
		return err
	}
	for i, tag := range tags {
		tagWithIndex := &model.TagWithIndex{
			Tag:   tag,
			Index: i,
		}
		if _, exist, err := a.tagRepository.Get(ws, tag.ID); err != nil {
			return fmt.Errorf("failed to get tag. id:%d : %w", tag.ID, err)
		} else if exist {
			if err := a.tagRepository.Update(ws, tagWithIndex); err != nil {
				return fmt.Errorf("failed to set tags: %w", err)
			}
		} else {
			if _, err := a.tagRepository.Add(ws, tagWithIndex); err != nil {
				return fmt.Errorf("failed to set tags: %w", err)
			}
		}
	}
	return nil
}
