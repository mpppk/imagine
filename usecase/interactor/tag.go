package interactor

import (
	"fmt"
	"sort"

	"github.com/mpppk/imagine/domain/client"

	"github.com/mpppk/imagine/domain/query"
	"github.com/mpppk/imagine/domain/service/tagsvc"

	"github.com/mpppk/imagine/domain/model"
	"github.com/mpppk/imagine/domain/repository"
)

type Tag struct {
	// FIXME: use client
	assetRepository repository.Asset
	tagRepository   repository.Tag
	tagQuery        query.Tag
}

func NewTag(c *client.Tag, asset repository.Asset) *Tag {
	// FIXME: use client
	return &Tag{
		assetRepository: asset,
		tagRepository:   c.TagRepository,
		tagQuery:        c.TagQuery,
	}
}

func (a *Tag) List(ws model.WSName) (tags []*model.Tag, err error) {
	tagsWithIndex, err := a.tagQuery.ListAll(ws)
	if err != nil {
		return nil, err
	}

	sort.Slice(tagsWithIndex, func(i, j int) bool {
		return tagsWithIndex[i].Index < tagsWithIndex[j].Index
	})

	tags = append(tags, tagsWithIndex...)

	return
}

// SaveTags persists tags and return saved tags.
// For each tag, if the tag already exists, update it. Otherwise add new tag with new ID.
// If model.Tag has non zero value but the tag which have the ID is not persisted, fail immediately and remain tags are not proceeded.
func (a *Tag) SaveTags(ws model.WSName, tags []*model.Tag) (newTags []*model.Tag, err error) {
	errMsg := "failed to save tags"
	for _, tag := range tags {
		if _, exist, err := a.tagQuery.Get(ws, tag.ID); err != nil {
			return nil, fmt.Errorf("failed to get tag. id:%d : %w", tag.ID, err)
		} else if exist {
			newTag, err := a.tagRepository.Save(ws, tag)
			if err != nil {
				return nil, fmt.Errorf("failed to set tags: %w", err)
			}
			newTags = append(newTags, newTag)
		} else {
			unregisteredTag, err := tag.SafeUnregister()
			if err != nil {
				return nil, fmt.Errorf("%s: %w", errMsg, err)
			}
			newTag, err := a.tagRepository.AddWithIndex(ws, unregisteredTag)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", errMsg, err)
			}
			newTags = append(newTags, newTag)
		}
	}
	return
}

// SetTags persists tags and return saved tags.
func (a *Tag) SetTags(ws model.WSName, tags []*model.UnindexedTag) (newTags []*model.Tag, err error) {
	errMsg := "failed to set tags"
	tagSet, err := a.tagQuery.ListAsSet(ws)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}

	tags2 := tagsvc.ToTagsWithOrderIndex(tags)
	_, nonExistsTagSet := tagSet.SplitByID(tagsvc.ToTagIDList(tags2))

	newTags, err = a.SaveTags(ws, tags2)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}

	// delete tag which does not have given name
	if err := a.tagRepository.Delete(ws, tagsvc.ToTagIDList(nonExistsTagSet.ToTags())); err != nil {
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}

	return
}

// SetTagByNames set provided tags.
// All existing tags will be replaced. (Internally, this method recreate tag bucket)
func (a *Tag) SetTagByNames(ws model.WSName, tagNames []string) (tags []*model.Tag, err error) {
	errMsg := "failed to set tags"
	tagSet, err := a.tagQuery.ListAsSet(ws)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}

	var tags2 []*model.UnindexedTag
	for _, tagName := range tagNames {
		tag, err := model.NewUnindexedTag(0, tagName)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", errMsg, err)
		}
		if t, ok := tagSet.GetByName(tagName); ok {
			tag.ID = t.ID
		}
		tags2 = append(tags2, tag)
	}

	tags, err = a.SetTags(ws, tags2)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}
	return
}

// Delete deletes and returns tags which match queries
func (a *Tag) Delete(ws model.WSName, queries []*model.Query) (deletedTags []*model.Tag, err error) {
	errMsg := "failed to delete tags"
	tags, err := a.tagQuery.ListByQueries(ws, queries)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}

	if err := a.assetRepository.UnAssignTags(ws, tagsvc.ToTagIDList(tags)...); err != nil {
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}

	if err := a.tagRepository.Delete(ws, tagsvc.ToTagIDList(tags)); err != nil {
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}

	return tags, nil
}
