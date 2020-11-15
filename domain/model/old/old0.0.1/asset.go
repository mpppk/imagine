package old0_0_1

import (
	"encoding/json"
	"fmt"

	"github.com/mpppk/imagine/domain/model"
)

type BoundingBoxID uint64
type BoundingBox struct {
	ID     BoundingBoxID `json:"id"`
	Tag    *Tag          `json:"tag"`
	X      int           `json:"x"`
	Y      int           `json:"y"`
	Width  int           `json:"width"`
	Height int           `json:"height"`
}

type TagID uint64
type Tag struct {
	ID   TagID  `json:"id"`
	Name string `json:"name"`
}

func (b *BoundingBox) Migrate() *model.BoundingBox {
	return &model.BoundingBox{
		ID:     model.BoundingBoxID(b.ID),
		TagID:  model.TagID(b.Tag.ID),
		X:      b.X,
		Y:      b.Y,
		Width:  b.Width,
		Height: b.Height,
	}
}

type AssetID uint64

type Asset struct {
	ID            AssetID        `json:"id"`
	Name          string         `json:"name"`
	Path          string         `json:"path"`
	BoundingBoxes []*BoundingBox `json:"boundingBoxes"`
}

func NewAssetFromJson(contents []byte) (*Asset, error) {
	var asset Asset
	if err := json.Unmarshal(contents, &asset); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json to asset: %s: %w", string(contents), err)
	}
	return &asset, nil
}

func (a *Asset) Migrate() *model.Asset {
	var boxes []*model.BoundingBox
	for _, box := range a.BoundingBoxes {
		boxes = append(boxes, box.Migrate())
	}

	return &model.Asset{
		ID:            model.AssetID(a.ID),
		Name:          a.Name,
		Path:          a.Path,
		BoundingBoxes: boxes,
	}
}
