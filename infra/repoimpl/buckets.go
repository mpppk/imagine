package repoimpl

import "github.com/mpppk/imagine/domain/model"

const (
	assetBucketName = "Asset"
	tagBucketName   = "Tag"
)

func createAssetBucketNames(ws model.WSName) []string {
	return []string{string(ws), assetBucketName}
}

func createTagBucketNames(ws model.WSName) []string {
	return []string{string(ws), tagBucketName}
}
