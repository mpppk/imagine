package model

type Tag string
type AssetID uint64

type Asset struct {
	ID   AssetID `json:"id"`
	Name string  `json:"name"`
	Path string  `json:"path"`
	Tags []Tag   `json:"tags"`
}
