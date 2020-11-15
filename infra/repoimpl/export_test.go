package repoimpl

import "github.com/mpppk/imagine/domain/model"

func CreateBoundingBox(id int) *model.BoundingBox {
	return &model.BoundingBox{
		ID:    model.BoundingBoxID(id),
		TagID: model.TagID(id),
	}
}
