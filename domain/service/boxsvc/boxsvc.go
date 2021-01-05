package boxsvc

import (
	"github.com/mpppk/imagine/domain/model"
)

func ToUniqTagNames(boxes []*model.ImportBoundingBox) (tagNames []string) {
	m := map[string]struct{}{}
	for _, box := range boxes {
		if box.HasTagName() {
			m[box.TagName] = struct{}{}
		}
	}
	for tagName := range m {
		tagNames = append(tagNames, tagName)
	}
	return
}
