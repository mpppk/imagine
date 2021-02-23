package blt

import "github.com/mpppk/imagine/domain/model"

const (
	AssetBucketName      = "Asset"
	PathBucketName       = "Path"
	TagBucketName        = "Tag"
	TagHistoryBucketName = "TagHistory"
	MetaBucketName       = "Meta"
)

func CreateAssetBucketNames(ws model.WSName) []string {
	return []string{string(ws), AssetBucketName}
}

func CreatePathBucketNames(ws model.WSName) []string {
	return []string{string(ws), PathBucketName}
}

func CreateTagBucketNames(ws model.WSName) []string {
	return []string{string(ws), TagBucketName}
}

func CreateTagHistoryBucketNames(ws model.WSName) []string {
	return []string{string(ws), TagHistoryBucketName}
}

func CreateMetaBucketNames() []string {
	return []string{MetaBucketName}
}
