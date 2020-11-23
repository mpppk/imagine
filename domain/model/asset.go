package model

import (
	"path/filepath"
	"strings"
)

type BoundingBoxID uint64
type BoundingBox struct {
	ID     BoundingBoxID `json:"id"`
	TagID  TagID         `json:"tagID"`
	X      int           `json:"x"`
	Y      int           `json:"y"`
	Width  int           `json:"width"`
	Height int           `json:"height"`
}

type ImportBoundingBox struct {
	*BoundingBox
	TagName string
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

type TagWithIndex struct {
	*Tag
	Index int
}

type TagSet struct {
	m     map[TagID]*Tag
	nameM map[string]*Tag
}

func NewTagSet() *TagSet {
	return &TagSet{
		m:     map[TagID]*Tag{},
		nameM: map[string]*Tag{},
	}
}

func (t *TagSet) Set(tag *Tag) bool {
	if sameNamedTag, ok := t.nameM[tag.Name]; ok && sameNamedTag.ID != tag.ID {
		return false
	}
	t.m[tag.ID] = tag
	t.nameM[tag.Name] = tag
	return true
}

func (t *TagSet) Get(id TagID) (*Tag, bool) {
	tag, ok := t.m[id]
	return tag, ok
}

func (t *TagSet) GetByName(name string) (*Tag, bool) {
	tag, ok := t.nameM[name]
	return tag, ok
}

func (t *TagSet) SubSetBy(f func(tag *Tag) bool) *TagSet {
	subset := NewTagSet()
	for _, tag := range t.m {
		if f(tag) {
			subset.Set(tag)
		}
	}
	return subset
}

type AssetID uint64

func AssetIDListToUint64List(assetIDList []AssetID) (idList []uint64) {
	for _, id := range assetIDList {
		idList = append(idList, uint64(id))
	}
	return
}

type Asset struct {
	ID            AssetID        `json:"id,omitempty"`
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

func (a *Asset) HasTag(tagID TagID) bool {
	for _, box := range a.BoundingBoxes {
		if box.TagID == tagID {
			return true
		}
	}
	return false
}

func (a *Asset) HasAnyOneOfTagID(tagSet *TagSet) bool {
	for _, box := range a.BoundingBoxes {
		if _, ok := tagSet.Get(box.TagID); ok {
			return true
		}
	}
	return false
}

type ImportAsset struct {
	*Asset        `mapstructure:",squash"`
	BoundingBoxes []*ImportBoundingBox `json:"boundingBoxes"`
}

func (a *ImportAsset) ToAsset() *Asset {
	var boxes []*BoundingBox

	for _, box := range a.BoundingBoxes {
		boxes = append(boxes, box.BoundingBox)
	}

	return &Asset{
		ID:            a.ID,
		Name:          a.Name,
		Path:          a.Path,
		BoundingBoxes: boxes,
	}
}

func ReplaceBoundingBoxByID(boxes []*BoundingBox, replaceBox *BoundingBox) (newBoxes []*BoundingBox) {
	for _, box := range boxes {
		if box.ID == replaceBox.ID {
			newBoxes = append(newBoxes, replaceBox)
		} else {
			newBoxes = append(newBoxes, box)
		}
	}
	return
}

func RemoveBoundingBoxByID(boxes []*BoundingBox, replaceBoxID BoundingBoxID) (newBoxes []*BoundingBox) {
	for _, box := range boxes {
		if box.ID != replaceBoxID {
			newBoxes = append(newBoxes, box)
		}
	}
	return
}

func NewAssetFromFilePath(filePath string) *Asset {
	name := strings.Replace(filepath.Base(filePath), filepath.Ext(filePath), "", -1)
	return &Asset{
		Name: name,
		Path: filePath,
	}
}
