package model

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
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

func NewTagSet(tags []*Tag) *TagSet {
	tagSet := &TagSet{
		m:     map[TagID]*Tag{},
		nameM: map[string]*Tag{},
	}
	for _, tag := range tags {
		tagSet.Set(tag)
	}
	return tagSet
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
	subset := NewTagSet(nil)
	for _, tag := range t.m {
		if f(tag) {
			subset.Set(tag)
		}
	}
	return subset
}

func (t *TagSet) ToMap() (map[TagID]*Tag, map[string]*Tag) {
	return t.m, t.nameM
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

func (a *Asset) ToJson() (string, error) {
	contents, err := json.Marshal(a)
	if err != nil {
		return "", fmt.Errorf("failed to marshal asset to json: %w", err)
	}
	return string(contents), nil
}

func (a *Asset) ToCSVRow(tagSet *TagSet) (string, error) {
	var tagNames []string
	for _, tagID := range BoxesToTagIDList(a.BoundingBoxes) {
		tag, ok := tagSet.Get(tagID)
		if !ok {
			log.Printf("warning: tag not found. id:%v", tagID)
			continue
		}
		tagNames = append(tagNames, tag.Name)
	}

	line := []string{
		strconv.Quote(strconv.Itoa(int(a.ID))),
		strconv.Quote(a.Path),
		strconv.Quote(strings.Join(tagNames, ",")),
	}

	return strings.Join(line, ","), nil
}

type ImportAsset struct {
	*Asset        `mapstructure:",squash"`
	BoundingBoxes []*ImportBoundingBox `json:"boundingBoxes"`
}

func (a *ImportAsset) ToAsset(tagSet *TagSet) (*Asset, error) {
	var boxes []*BoundingBox

	for _, box := range a.BoundingBoxes {
		newBox := box.BoundingBox
		if newBox == nil {
			newBox = &BoundingBox{}
		}
		if newBox.TagID == 0 {
			tag, ok := tagSet.GetByName(box.TagName)
			if !ok {
				return nil, fmt.Errorf("unknown tag name(%s)", box.TagName)
			}
			newBox.TagID = tag.ID
		}
		boxes = append(boxes, newBox)
	}

	return &Asset{
		ID:            a.ID,
		Name:          a.Name,
		Path:          a.Path,
		BoundingBoxes: boxes,
	}, nil
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

func NewImportAssetFromFilePath(filePath string) *ImportAsset {
	name := strings.Replace(filepath.Base(filePath), filepath.Ext(filePath), "", -1)
	return &ImportAsset{
		Asset: &Asset{
			Name: name,
			Path: filePath,
		},
	}
}

func BoxesToTagIDList(boxes []*BoundingBox) (idList []TagID) {
	tagM := map[TagID]struct{}{}
	for _, box := range boxes {
		tagM[box.TagID] = struct{}{}
	}

	for id := range tagM {
		idList = append(idList, id)
	}
	return
}
