package tagsvc

import "github.com/mpppk/imagine/domain/model"

func ToTagNames(tags []*model.Tag) (tagNames []string) {
	for _, tag := range tags {
		tagNames = append(tagNames, tag.Name)
	}
	return
}

func ToTagIDList(tags []*model.TagWithIndex) (idList []model.TagID) {
	for _, tag := range tags {
		idList = append(idList, tag.ID)
	}
	return
}
