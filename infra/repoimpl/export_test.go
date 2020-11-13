package repoimpl

import "github.com/mpppk/imagine/domain/model"

func CreateBoundingBox(id int, tagName string) *model.BoundingBox {
	return &model.BoundingBox{
		ID: model.BoundingBoxID(id),
		Tag: &model.Tag{
			ID:   model.TagID(id),
			Name: tagName,
		},
	}
}
