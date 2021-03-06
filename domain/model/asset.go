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

func (b *BoundingBox) HasID() bool {
	return b.ID != 0
}

func (b *BoundingBox) HasTagID() bool {
	return b.TagID != 0
}

// IsSame checks if tow boxes are the same.
// This method ignore ID.
func (b *BoundingBox) IsSame(box *BoundingBox) bool {
	return b.TagID == box.TagID &&
		b.X == box.X && b.Y == box.Y &&
		b.Width == box.Width && b.Height == box.Height
}

type BoundingBoxes []*BoundingBox

// Merge merge boxes
func (boxes BoundingBoxes) Merge(otherBoxes []*BoundingBox) []*BoundingBox {
	var newBoxes BoundingBoxes = make([]*BoundingBox, len(boxes))
	copy(newBoxes, boxes)
	for _, box := range otherBoxes {
		if !newBoxes.HasSameBox(box) {
			newBoxes = append(newBoxes, box)
		}
	}
	return newBoxes
}

func (boxes BoundingBoxes) HasSameBox(box *BoundingBox) bool {
	for _, boundingBox := range boxes {
		if boundingBox.IsSame(box) {
			return true
		}
	}
	return false
}

type ImportBoundingBox struct {
	*BoundingBox `mapstructure:",squash"`
	TagName      string `json:"tagName"`
}

func NewImportBoundingBoxFromTagID(tagID TagID) *ImportBoundingBox {
	return &ImportBoundingBox{
		BoundingBox: &BoundingBox{
			TagID: tagID,
		},
	}
}

func (b *ImportBoundingBox) Validate(tagSet *TagSet) error {
	if !b.HasTagName() && !b.HasTagID() {
		return fmt.Errorf("bouding box's tag name and tag id are empty")
	}

	if tagSet != nil {
		if b.HasTagName() {
			if _, ok := tagSet.GetByName(b.TagName); !ok {
				return fmt.Errorf("tag name(%s) not found in tag set", b.TagName)
			}
		}

		if b.HasTagID() {
			if _, ok := tagSet.Get(b.TagID); !ok {
				return fmt.Errorf("tag id(%d) not found in tag set", b.TagID)
			}
		}

		if b.HasTagName() && b.HasTagID() {
			if tag, _ := tagSet.Get(b.TagID); tag.Name != b.TagName {
				return fmt.Errorf("tag id and name inconsistency. provided id:%d, name:%s, but stored name is %s", b.TagID, b.TagName, tag.Name)
			}
		}
	}

	return nil
}

func (b *ImportBoundingBox) HasTagName() bool {
	return b.TagName != ""
}

type AssetID uint64

func NewAssetID(id int) (AssetID, error) {
	if id < 0 {
		return 0, fmt.Errorf("negative number can not be AssetID: %d", id)
	}
	if id == 0 {
		return 0, fmt.Errorf("zero can not be AssetID")
	}
	return AssetID(id), nil
}

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

func NewAssetFromBytes(bytes []byte) (*Asset, error) {
	var asset Asset
	if err := json.Unmarshal(bytes, &asset); err != nil {
		return nil, fmt.Errorf("failed to unmarshal asset: %w", err)
	}
	return &asset, nil
}

func (a *Asset) Validate() error {
	if a.ID == 0 && a.Path == "" {
		return fmt.Errorf("id and path are empty")
	}
	return nil
}

func (a *Asset) IsAddable() bool {
	return a != nil && !a.HasID() && a.HasPath()
}

// IsUpdatableByID checks if this asset can be updated.
// If asset or box which the asset has does not have ID, the asset is not updatable.
func (a *Asset) IsUpdatableByID() bool {
	if a == nil || !a.HasID() {
		return false
	}

	for _, box := range a.BoundingBoxes {
		if !box.HasTagID() {
			return false
		}
	}
	return true
}

func (a *Asset) IsSavable() bool {
	if !a.IsUpdatableByID() {
		return false
	}
	return a.HasPath()
}

func (a *Asset) GetID() uint64 {
	return uint64(a.ID)
}

func (a *Asset) HasID() bool {
	return a.ID != 0
}

func (a *Asset) SetID(id uint64) {
	a.ID = AssetID(id)
}

func (a *Asset) HasPath() bool {
	return a.Path != ""
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

func (a *Asset) CanBeUpdatedBy(asset *Asset) (ok bool, reason string) {
	if asset.HasID() && a.ID != asset.ID {
		return false, fmt.Sprintf("IDs are different(%d, %d)", a.ID, asset.ID)
	}
	return true, ""
}

// UpdateBy merge the itself and argument asset properties. This is destructive method.
// If receiver or arg asset is nil, UpdateBy method do nothing.
func (a *Asset) UpdateBy(asset *Asset) error {
	errMsg := "failed to update asset"
	if a == nil || asset == nil {
		return nil
	}

	if ok, reason := a.CanBeUpdatedBy(asset); !ok {
		return fmt.Errorf("%s: %s", errMsg, reason)
	}

	if asset.HasPath() {
		a.Path = asset.Path
		a.Name = assetPathToName(a.Path)
	}

	if asset.BoundingBoxes != nil {
		if a.BoundingBoxes == nil {
			a.BoundingBoxes = asset.BoundingBoxes
		} else {
			a.BoundingBoxes = (BoundingBoxes)(a.BoundingBoxes).Merge(asset.BoundingBoxes)
		}
	}
	return nil
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

func assetPathToName(p string) string {
	return strings.Replace(filepath.Base(p), filepath.Ext(p), "", -1)
}

type ImportAsset struct {
	*Asset        `mapstructure:",squash"`
	BoundingBoxes []*ImportBoundingBox `json:"boundingBoxes"`
}

func NewImportAsset(id AssetID, path string, boxes []*ImportBoundingBox) *ImportAsset {
	a := NewImportAssetFromFilePath(path)
	if id != 0 {
		a.ID = id
	}
	a.BoundingBoxes = boxes
	return a
}

func NewImportAssetFromJson(contents []byte) (*ImportAsset, error) {
	var asset ImportAsset
	if err := json.Unmarshal(contents, &asset); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json to import asset")
	}

	return &asset, asset.Validate(nil)
}

func (a *ImportAsset) Validate(tagSet *TagSet) error {
	asset := a.Asset
	if asset == nil {
		asset = &Asset{}
	}
	if err := asset.Validate(); err != nil {
		return err
	}

	for _, box := range a.BoundingBoxes {
		if err := box.Validate(tagSet); err != nil {
			return err
		}
	}

	return nil
}

func (a *ImportAsset) ToAsset(tagSet *TagSet) (*Asset, error) {
	var boxes []*BoundingBox

	if a.Asset.BoundingBoxes != nil {
		return nil, fmt.Errorf("failed to convert import asset to asset. import asset should not have bounding box, instead use import bounding box: %#v", a.Asset.BoundingBoxes)
	}

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
	name := assetPathToName(filePath)
	return &Asset{
		Name: name,
		Path: filePath,
	}
}

func NewImportAssetFromFilePath(filePath string) *ImportAsset {
	name := assetPathToName(filePath)
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
