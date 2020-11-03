package repoimpl

import "github.com/mpppk/imagine/domain/model"

const (
	assetBucketName = "Asset"
	pathBucketName  = "Path"
	tagBucketName   = "Tag"
)

func createAssetBucketNames(ws model.WSName) []string {
	return []string{string(ws), assetBucketName}
}

func createPathBucketNames(ws model.WSName) []string {
	return []string{string(ws), pathBucketName}
}

func createTagBucketNames(ws model.WSName) []string {
	return []string{string(ws), tagBucketName}
}
