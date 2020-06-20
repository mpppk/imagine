package model

type Tag string
type AssetID uint64

type Asset struct {
	ID   AssetID `boltholdKey:"ID"`
	Name string
	Path string
	Tags []Tag
}
