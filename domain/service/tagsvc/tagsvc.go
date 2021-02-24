package tagsvc

import "github.com/mpppk/imagine/domain/model"

func ToTagNames(tags []*model.Tag) (tagNames []string) {
	for _, tag := range tags {
		tagNames = append(tagNames, tag.Name)
	}
	return
}

func ToTagsWithOrderIndex(unIndexedTags []*model.UnindexedTag) (tags []*model.Tag) {
	for i, tag := range unIndexedTags {
		// error can be ignored because i always be positive
		tag, _ := tag.Index(i)
		tags = append(tags, tag)
	}
	return
}

func ToTagIDList(tags []*model.Tag) (idList []model.TagID) {
	for _, tag := range tags {
		idList = append(idList, tag.ID)
	}
	return
}
