package repoimpl

import (
	"encoding/json"
	"fmt"

	"github.com/mpppk/imagine/domain/model"
	"github.com/mpppk/imagine/domain/repository"
	bolt "go.etcd.io/bbolt"
)

const globalBucketName = "WorkSpace"

var workSpacesKey = []byte("WorkSpaces")

var globalBucketNames = []string{globalBucketName}

type BBoltWorkSpace struct {
	base           *boltRepository
	pathRepository *pathRepository
}

func NewBBoltWorkSpace(b *bolt.DB) repository.WorkSpace {
	return &BBoltWorkSpace{
		base:           newBoltRepository(b),
		pathRepository: newPathRepository(b),
	}
}

func (b *BBoltWorkSpace) globalBucketFunc(f func(bucket *bolt.Bucket) error) error {
	return b.base.bucketFunc(globalBucketNames, f)
}

func (b *BBoltWorkSpace) getWorkSpacesFromBucket(bucket *bolt.Bucket) (workspaces []*model.WorkSpace, err error) {
	workSpacesBytes := bucket.Get(workSpacesKey)
	if workSpacesBytes == nil {
		return nil, nil
	}
	err = json.Unmarshal(workSpacesBytes, &workspaces)
	return
}

func (b *BBoltWorkSpace) setWorkSpaces(bucket *bolt.Bucket, workspaces []*model.WorkSpace) error {
	workspacesBytes, err := json.Marshal(workspaces)
	if err != nil {
		return fmt.Errorf("failed to marshal workspaces: %w", err)
	}
	return bucket.Put(workSpacesKey, workspacesBytes)
}

func (b *BBoltWorkSpace) updateWorkSpaces(f func(workspaces []*model.WorkSpace) ([]*model.WorkSpace, error)) error {
	return b.globalBucketFunc(func(bucket *bolt.Bucket) error {
		workspaces, err := b.getWorkSpacesFromBucket(bucket)
		if err != nil {
			return err
		}
		newWorkspaces, err := f(workspaces)
		if err != nil {
			return err
		}
		return b.setWorkSpaces(bucket, newWorkspaces)
	})
}

func (b *BBoltWorkSpace) List() (workspaces []*model.WorkSpace, err error) {
	err = b.globalBucketFunc(func(bucket *bolt.Bucket) error {
		workspaces, err = b.getWorkSpacesFromBucket(bucket)
		return err
	})
	return
}

func (b *BBoltWorkSpace) Add(ws *model.WorkSpace) error {
	return b.updateWorkSpaces(func(workspaces []*model.WorkSpace) ([]*model.WorkSpace, error) {
		return append(workspaces, ws), nil
	})
}

func (b *BBoltWorkSpace) Update(ws *model.WorkSpace) error {
	return b.updateWorkSpaces(func(workspaces []*model.WorkSpace) ([]*model.WorkSpace, error) {
		return replaceWorkSpaceByName(workspaces, ws), nil
	})
}

func replaceWorkSpaceByName(workspaces []*model.WorkSpace, ws *model.WorkSpace) (newWorkspaces []*model.WorkSpace) {
	for _, workspace := range workspaces {
		if workspace.Name == ws.Name {
			newWorkspaces = append(newWorkspaces, ws)
		} else {
			newWorkspaces = append(newWorkspaces, workspace)
		}
	}
	return
}
