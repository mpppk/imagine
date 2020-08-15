package model

type BoundingBoxID uint64
type BoundingBox struct {
	ID     BoundingBoxID
	Tag    *Tag
	X      int
	Y      int
	Width  int
	Height int
}

type TagID uint64
type Tag struct {
	ID   TagID
	Name string
}

func (t *Tag) GetID() uint64 {
	return uint64(t.ID)
}

func (t *Tag) SetID(id uint64) {
	t.ID = TagID(id)
}

type AssetID uint64

type Asset struct {
	ID            AssetID        `json:"id"`
	Name          string         `json:"name"`
	Path          string         `json:"path"`
	BoundingBoxes []*BoundingBox `json:"boundingBoxes"`
}
