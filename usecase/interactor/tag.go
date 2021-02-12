package interactor

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

// PutTags add or update tags.
// If tag which have specified id is exist, it will be updated. Otherwise add new tag.
func (a *Tag) PutTags(ws model.WSName, tags []*model.Tag) error {
	for i, tag := range tags {
		tagWithIndex := &model.TagWithIndex{
			Tag:   tag,
			Index: i,
		}
		if _, exist, err := a.tagRepository.Get(ws, tag.ID); err != nil {
			return fmt.Errorf("failed to get tag. id:%d : %w", tag.ID, err)
		} else if exist {
			if err := a.tagRepository.Put(ws, tagWithIndex); err != nil {
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

// SetTags set provided tags.
// All existing tags will be replaced. (Internally, this method recreate tag bucket)
func (a *Tag) SetTags(ws model.WSName, tagNames []string) ([]model.TagID, error) {
	if err := a.tagRepository.RecreateBucket(ws); err != nil {
		return nil, fmt.Errorf("failed to recreate tag bucket. ws: %s", ws)
	}
	idList, err := a.tagRepository.AddByNames(ws, tagNames)
	if err != nil {
		return nil, fmt.Errorf("failed to add tags by name. ws: %s, tag names: %v", ws, tagNames)
	}
	return idList, nil
}
