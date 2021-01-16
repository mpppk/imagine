package model

type WSName string

func (w WSName) IsValid() bool {
	return w != ""
}

type WorkSpaceID uint64
type WorkSpace struct {
	ID       WorkSpaceID `json:"id"`
	Name     WSName      `json:"name"`
	BasePath string      `json:"basePath"`
}
