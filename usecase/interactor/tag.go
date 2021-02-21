package interactor

import (
	"fmt"
	"sort"

	"github.com/mpppk/imagine/domain/service/assetsvc/tagsvc"

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
func (a *Tag) SetTags(ws model.WSName, tagNames []string) (idList []model.TagID, err error) {
	errMsg := "failed to set tags"
	//a.updateByTagSetFunc(ws, func(tagSet *model.TagSet) (*model.TagSet, error) {
	//	newTagSet := model.NewTagSet(nil)
	//	for _, tagName := range tagNames {
	//		newTagSet.Set()
	//	}
	//	return tagSet, nil
	//})
	tagSet, err := a.tagRepository.ListAsSet(ws)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}

	existsTagSet, nonExistsTagSet := tagSet.SplitByNames(tagNames)

	var tagNamesToBeAdded []string
	for _, tagName := range tagNames {
		if tag, ok := existsTagSet.GetByName(tagName); ok {
			idList = append(idList, tag.ID)
			// FIXME: If the Tag has properties other than ID and Name in the future, it should be updated.
			continue
		}
		// add tag if it does not exist yet
		tagNamesToBeAdded = append(tagNamesToBeAdded, tagName)
	}

	idList2, err := a.tagRepository.AddByNames(ws, tagNamesToBeAdded)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}

	idList = append(idList, idList2...)

	// delete tags which does not have given name
	if err := a.tagRepository.Delete(ws, tagsvc.ToTagIDList(nonExistsTagSet.ToTags())); err != nil {
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}

	return
}

//func (a *Tag) updateByTagSetFunc(ws model.WSName, f func(tagSet *model.TagSet) (*model.TagSet, error)) error {
//	errMsg := "updateByTagSetFunc failed"
//	tagSet, err := a.tagRepository.ListAsSet(ws)
//	if err != nil {
//		return fmt.Errorf("%s: %w", errMsg, err)
//	}
//
//	newTagSet, err := f(tagSet)
//	if err != nil {
//		return fmt.Errorf("%s: %w", errMsg, err)
//	}
//	return a.tagRepository.UpdateByTagSet(ws, newTagSet)
//}
