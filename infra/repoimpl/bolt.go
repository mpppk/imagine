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
func (b *boltRepository) bucket(tx *bolt.Tx) *bolt.Bucket {
	return tx.Bucket([]byte(assetBucketName))
}

func (b *boltRepository) bucketFunc(f func(bucket *bolt.Bucket) error) error {
	return b.bolt.Update(func(tx *bolt.Tx) error {
		return f(b.bucket(tx))
	})
}

func (b *boltRepository) loBucketFunc(f func(bucket *bolt.Bucket) error) error {
	return b.bolt.View(func(tx *bolt.Tx) error {
		return f(b.bucket(tx))
	})
}

func (b *boltRepository) createBucketIfNotExist() error {
	return b.bolt.Update(func(tx *bolt.Tx) error {
		if bucket := tx.Bucket([]byte(assetBucketName)); bucket != nil {
			return nil
		}
		_, err := tx.CreateBucket([]byte(assetBucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}

func (b *boltRepository) close() error {
	return b.bolt.Close()
}
