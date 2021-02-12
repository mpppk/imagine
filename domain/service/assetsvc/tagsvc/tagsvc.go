package tagsvc

import "github.com/mpppk/imagine/domain/model"

func ToTagNames(tags []*model.Tag) (tagNames []string) {
	for _, tag := range tags {
		tagNames = append(tagNames, tag.Name)
	}
	return
}
