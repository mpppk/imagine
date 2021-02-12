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

// SaveTags persists tags and return ID list of tags.
// For each tag, if the tag already exists, update it. Otherwise add new tag with new ID.
// If model.Tag has non zero value but the tag which have the ID is not persisted, fail immediately and remain tags are not proceeded.
func (a *Tag) SaveTags(ws model.WSName, tags []*model.Tag) (idList []model.TagID, err error) {
	for i, tag := range tags {
		tagWithIndex := &model.TagWithIndex{
			Tag:   tag,
			Index: i,
		}
		if _, exist, err := a.tagRepository.Get(ws, tag.ID); err != nil {
			return nil, fmt.Errorf("failed to get tag. id:%d : %w", tag.ID, err)
		} else if exist {
			id, err := a.tagRepository.Save(ws, tagWithIndex)
			if err != nil {
				return nil, fmt.Errorf("failed to set tags: %w", err)
			}
			idList = append(idList, id)
		} else {
			id, err := a.tagRepository.Add(ws, tagWithIndex)
			if err != nil {
				return nil, fmt.Errorf("failed to set tags: %w", err)
			}
			idList = append(idList, id)
		}
	}
	return
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
