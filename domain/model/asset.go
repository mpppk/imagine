package model

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

func (a *Asset) GetID() uint64 {
	return uint64(a.ID)
}

func (a *Asset) SetID(id uint64) {
	a.ID = AssetID(id)
}
