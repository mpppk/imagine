package repoimpl

import (
	"encoding/binary"
	"errors"
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

// bucket gets or creates buckets
func (b *boltRepository) bucket(bucketNames []string, tx *bolt.Tx) (*bolt.Bucket, error) {
	if len(bucketNames) == 0 {
		return nil, errors.New("empty bucket name is provided")
	}

	bucket, err := tx.CreateBucketIfNotExists([]byte(bucketNames[0]))
	if err != nil {
		return nil, fmt.Errorf("failed to create bucket(name: %s): %w", bucketNames[0], err)
	}

	if len(bucketNames) == 1 {
		return bucket, nil
	}

	for _, bucketName := range bucketNames[1:] {
		b, err := bucket.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket(name: %s): %w", bucketName, err)
		}
		bucket = b
	}
	return bucket, nil
}

func (b *boltRepository) internalBucketFunc(bucketNames []string, f func(bucket *bolt.Bucket) error) func(tx *bolt.Tx) error {
	return func(tx *bolt.Tx) error {
		bucket, err := b.bucket(bucketNames, tx)
		if err != nil {
			return fmt.Errorf("failed to get bucket from %v: %w", bucketNames, err)
		}
		return f(bucket)
	}
}
func (b *boltRepository) bucketFunc(bucketNames []string, f func(bucket *bolt.Bucket) error) error {
	return b.bolt.Update(b.internalBucketFunc(bucketNames, f))
}

func (b *boltRepository) loBucketFunc(bucketNames []string, f func(bucket *bolt.Bucket) error) error {
	return b.bolt.View(b.internalBucketFunc(bucketNames, f))
}

func (b *boltRepository) close() error {
	return b.bolt.Close()
}
