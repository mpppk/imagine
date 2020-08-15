package repoimpl

import (
	"encoding/binary"
	"encoding/json"
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

func (b *boltRepository) createBucketIfNotExist(bucketNames []string) error {
	return b.bolt.Update(func(tx *bolt.Tx) error {
		_, err := b.getOrCreateBucket(bucketNames, tx)
		return err
	})
}

// getOrCreateBucket gets or creates buckets
func (b *boltRepository) getOrCreateBucket(bucketNames []string, tx *bolt.Tx) (*bolt.Bucket, error) {
	if len(bucketNames) == 0 {
		return nil, errors.New("empty bucket name is provided")
	}

	bucket, err := tx.CreateBucketIfNotExists([]byte(bucketNames[0]))
	if err != nil {
		return nil, fmt.Errorf("failed to create getOrCreateBucket(name: %s): %w", bucketNames[0], err)
	}

	if len(bucketNames) == 1 {
		return bucket, nil
	}

	for _, bucketName := range bucketNames[1:] {
		b, err := bucket.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return nil, fmt.Errorf("failed to create getOrCreateBucket(name: %s): %w", bucketName, err)
		}
		bucket = b
	}
	return bucket, nil
}

// getBucket gets bucket
func (b *boltRepository) getBucket(bucketNames []string, tx *bolt.Tx) (*bolt.Bucket, error) {
	if len(bucketNames) == 0 {
		return nil, errors.New("empty bucket name is provided")
	}

	bucket := tx.Bucket([]byte(bucketNames[0]))
	if bucket == nil {
		return nil, fmt.Errorf("failed to get bucket(name: %s)", bucketNames[0])
	}

	if len(bucketNames) == 1 {
		return bucket, nil
	}

	for _, bucketName := range bucketNames[1:] {
		b := bucket.Bucket([]byte(bucketName))
		if b == nil {
			return nil, fmt.Errorf("failed to get bucket(name: %s)", bucketName)
		}
		bucket = b
	}
	return bucket, nil
}

func (b *boltRepository) internalBucketFunc(bucketNames []string, f func(bucket *bolt.Bucket) error) func(tx *bolt.Tx) error {
	return func(tx *bolt.Tx) error {
		bucket, err := b.getOrCreateBucket(bucketNames, tx)
		if err != nil {
			return fmt.Errorf("failed to get or create bucket from %v: %w", bucketNames, err)
		}
		return f(bucket)
	}
}

func (b *boltRepository) internalLOBucketFunc(bucketNames []string, f func(bucket *bolt.Bucket) error) func(tx *bolt.Tx) error {
	return func(tx *bolt.Tx) error {
		bucket, err := b.getBucket(bucketNames, tx)
		if err != nil {
			return fmt.Errorf("failed to get getOrCreateBucket from %v: %w", bucketNames, err)
		}
		return f(bucket)
	}
}

func (b *boltRepository) bucketFunc(bucketNames []string, f func(bucket *bolt.Bucket) error) error {
	return b.bolt.Update(b.internalBucketFunc(bucketNames, f))
}

func (b *boltRepository) loBucketFunc(bucketNames []string, f func(bucket *bolt.Bucket) error) error {
	return b.bolt.View(b.internalLOBucketFunc(bucketNames, f))
}

func (b *boltRepository) close() error {
	return b.bolt.Close()
}

type boltData interface {
	GetID() uint64
	SetID(id uint64)
}

func (b *boltRepository) add(bucketNames []string, data boltData) error {
	return b.bucketFunc(bucketNames, func(bucket *bolt.Bucket) error {
		id, err := bucket.NextSequence()
		if err != nil {
			return err
		}
		data.SetID(id)
		s, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to marshal data to json: %w", err)
		}
		return bucket.Put(itob(data.GetID()), s)
	})
}

func (b *boltRepository) get(bucketNames []string, id uint64) (data []byte, err error) {
	err = b.loBucketFunc(bucketNames, func(bucket *bolt.Bucket) error {
		data = bucket.Get(itob(id))
		if data == nil {
			return fmt.Errorf("data does not exist: %v", id)
		}
		return nil
	})
	return
}

func (b *boltRepository) update(bucketNames []string, data boltData) error {
	return b.bucketFunc(bucketNames, func(bucket *bolt.Bucket) error {
		s, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to marshal tag to json: %w", err)
		}
		return bucket.Put(itob(data.GetID()), s)
	})
}

func (b *boltRepository) recreateBucket(bucketNames []string) error {
	if len(bucketNames) == 0 {
		return fmt.Errorf("empty bucket name is given to deleteBucket")
	}

	return b.bucketFunc(bucketNames[:len(bucketNames)-1], func(bucket *bolt.Bucket) error {
		lastBucketName := bucketNames[len(bucketNames)-1]
		lastBucketNameBytes := []byte(lastBucketName)
		if err := bucket.DeleteBucket(lastBucketNameBytes); err != nil {
			return fmt.Errorf("failed to delete bucket. name: " + lastBucketName)
		}
		_, err := bucket.CreateBucket(lastBucketNameBytes)
		if err != nil {
			return fmt.Errorf("failed to recreate bucket. name: " + lastBucketName)
		}
		return nil
	})
}

// TODO: 次はこれを使ってtagを滅ぼし、再度登録するsetTagsをusecaseとして実装
