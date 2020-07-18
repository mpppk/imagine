package repoimpl

import (
	"encoding/binary"
	"fmt"

	bolt "go.etcd.io/bbolt"
)

// itob returns an 8-byte big endian representation of v.
func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

type boltRepository struct {
	bolt *bolt.DB
}

func newBoltRepository(b *bolt.DB) *boltRepository {
	return &boltRepository{
		bolt: b,
	}
}

func (b *boltRepository) bucket(bucketNames []string, tx *bolt.Tx) *bolt.Bucket {
	var bucket *bolt.Bucket = nil
	for _, bucketName := range bucketNames {
		bucket = tx.Bucket([]byte(bucketName))
	}
	return bucket
}

func (b *boltRepository) hasBucket(bucketNames []string, tx *bolt.Tx) bool {
	var bucket *bolt.Bucket = nil
	for _, bucketName := range bucketNames {
		bucket = tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return false
		}
	}
	return true
}

func (b *boltRepository) bucketFunc(bucketNames []string, f func(bucket *bolt.Bucket) error) error {
	return b.bolt.Update(func(tx *bolt.Tx) error {
		return f(b.bucket(bucketNames, tx))
	})
}

func (b *boltRepository) loBucketFunc(bucketNames []string, f func(bucket *bolt.Bucket) error) error {
	return b.bolt.View(func(tx *bolt.Tx) error {
		return f(b.bucket(bucketNames, tx))
	})
}

func (b *boltRepository) createBucketIfNotExist(bucketNames []string) error {
	return b.bolt.Update(func(tx *bolt.Tx) error {
		var bucket *bolt.Bucket = nil
		for _, bucketName := range bucketNames {
			bucket = tx.Bucket([]byte(bucketName))
			if bucket == nil {
				if _, err := tx.CreateBucket([]byte(bucketName)); err != nil {
					return fmt.Errorf("failed to create bucket: %s", err)
				}
			}
		}

		return nil
	})
}

func (b *boltRepository) close() error {
	return b.bolt.Close()
}
