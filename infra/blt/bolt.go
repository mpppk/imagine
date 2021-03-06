package blt

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"

	bolt "go.etcd.io/bbolt"
)

type BoltData interface {
	GetID() uint64
	SetID(id uint64)
}

type Repository struct {
	bolt *bolt.DB
}

func NewRepository(b *bolt.DB) *Repository {
	return &Repository{
		bolt: b,
	}
}

// ---- utilities ----

func Itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

func Btoi(bytes []byte) uint64 {
	padding := make([]byte, 8-len(bytes))
	i := binary.BigEndian.Uint64(append(padding, bytes...))
	return i
}

func toJson(data BoltData) ([]byte, error) {
	s, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal bolt data to json: %w", err)
	}
	return s, nil
}

func nextSequence(bucket *bolt.Bucket) (uint64, error) {
	id, err := bucket.NextSequence()
	if err != nil {
		return 0, err
	}
	if id == 0 {
		id, err = bucket.NextSequence()
		if err != nil {
			return 0, err
		}
	}
	return id, nil
}

func add(bucket *bolt.Bucket, data interface{}) (uint64, error) {
	id, err := nextSequence(bucket)
	if err != nil {
		return 0, err
	}
	s, err := json.Marshal(data)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal bolt data to json: %w", err)
	}
	if err := bucket.Put(Itob(id), s); err != nil {
		return 0, fmt.Errorf("failed to add data to bolt. ID:%d", id)
	}
	return id, nil
}

// addWithID adds data to bolt. ID will be generated and set to data. So this method change ID of given data.
// Even if ID is given, it is always ignored.
func addWithID(bucket *bolt.Bucket, unregisteredData BoltData) (uint64, error) {
	id, err := nextSequence(bucket)
	if err != nil {
		return 0, err
	}
	unregisteredData.SetID(id)
	if err := putByID(bucket, unregisteredData); err != nil {
		return 0, fmt.Errorf("failed to ")
	}
	return id, nil
}

func putByID(bucket *bolt.Bucket, data BoltData) error {
	s, err := toJson(data)
	if err != nil {
		return err
	}
	if err := bucket.Put(Itob(data.GetID()), s); err != nil {
		return fmt.Errorf("failed to put data to bolt. ID:%d", data.GetID())
	}
	return nil
}

// ---- bucket operations ----

func (b *Repository) CreateBucketIfNotExist(bucketNames []string) error {
	return b.bolt.Update(func(tx *bolt.Tx) error {
		_, err := b.getOrCreateBucket(bucketNames, tx)
		return err
	})
}

// getOrCreateBucket gets or creates buckets
func (b *Repository) getOrCreateBucket(bucketNames []string, tx *bolt.Tx) (*bolt.Bucket, error) {
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
func (b *Repository) getBucket(bucketNames []string, tx *bolt.Tx) (*bolt.Bucket, error) {
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

func (b *Repository) internalBucketFunc(bucketNames []string, f func(bucket *bolt.Bucket) error) func(tx *bolt.Tx) error {
	return func(tx *bolt.Tx) error {
		bucket, err := b.getOrCreateBucket(bucketNames, tx)
		if err != nil {
			return fmt.Errorf("failed to get or create bucket from %v: %w", bucketNames, err)
		}
		return f(bucket)
	}
}

func (b *Repository) internalLOBucketFunc(bucketNames []string, f func(bucket *bolt.Bucket) error) func(tx *bolt.Tx) error {
	return func(tx *bolt.Tx) error {
		bucket, err := b.getBucket(bucketNames, tx)
		if err != nil {
			return fmt.Errorf("failed to get getOrCreateBucket from %v: %w", bucketNames, err)
		}
		return f(bucket)
	}
}

func (b *Repository) BucketFunc(bucketNames []string, f func(bucket *bolt.Bucket) error) error {
	f2 := func(bucket *bolt.Bucket) error {
		e := f(bucket)
		return e
	}
	e := b.bolt.Update(b.internalBucketFunc(bucketNames, f2))
	return e
}

func (b *Repository) LoBucketFunc(bucketNames []string, f func(bucket *bolt.Bucket) error) error {
	f2 := func(bucket *bolt.Bucket) error {
		e := f(bucket)
		return e
	}
	return b.bolt.View(b.internalLOBucketFunc(bucketNames, f2))
}

func (b *Repository) RecreateBucket(bucketNames []string) error {
	if len(bucketNames) == 0 {
		return fmt.Errorf("empty bucket name is given to deleteBucket")
	}

	return b.BucketFunc(bucketNames[:len(bucketNames)-1], func(bucket *bolt.Bucket) error {
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

func (b *Repository) batchBucketFunc(bucketNames []string, f func(bucket *bolt.Bucket) error) error {
	f2 := func(bucket *bolt.Bucket) error {
		e := f(bucket)
		return e
	}
	e := b.bolt.Batch(b.internalBucketFunc(bucketNames, f2))
	return e
}

// ---- data operations ----

func (b *Repository) AddIntByString(bucketNames []string, k string, v uint64) error {
	return b.BucketFunc(bucketNames, func(bucket *bolt.Bucket) error {
		e := bucket.Put([]byte(k), Itob(v))
		return e
	})
}

// add adds data to bolt and set new ID which generated by bolt to data. So this method modifies data argument.
// This method always assign new ID to data, so even if already data have ID other than 0, it will be ignored and overwritten.
func (b *Repository) Add(bucketNames []string, data interface{}) (uint64, error) {
	var retId uint64
	e := b.BucketFunc(bucketNames, func(bucket *bolt.Bucket) error {
		id, err := add(bucket, data)
		retId = id
		return err
	})
	return retId, e
}

// addWithID adds data to bolt and set new ID which generated by bolt to data. So this method modifies data argument.
// This method always assign new ID to data, so even if already data have ID other than 0, it will be ignored and overwritten.
func (b *Repository) AddWithID(bucketNames []string, data BoltData) (uint64, error) {
	var retId uint64
	e := b.BucketFunc(bucketNames, func(bucket *bolt.Bucket) error {
		id, err := addWithID(bucket, data)
		retId = id
		return err
	})
	return retId, e
}

func (b *Repository) Get(bucketNames []string, id uint64) (data []byte, exist bool, err error) {
	err = b.LoBucketFunc(bucketNames, func(bucket *bolt.Bucket) error {
		data = bucket.Get(Itob(id))
		return nil
	})
	return data, data != nil, err
}

func (b *Repository) GetIDByString(bucketNames []string, key string) (id uint64, exist bool, err error) {
	var data []byte
	err = b.LoBucketFunc(bucketNames, func(bucket *bolt.Bucket) error {
		data = bucket.Get([]byte(key))
		return nil
	})
	return Btoi(data), data != nil, err
}

func (b *Repository) List(bucketNames []string) (dataList [][]byte, err error) {
	err = b.ForEach(bucketNames, func(value []byte) error {
		dataList = append(dataList, value)
		return nil
	})
	return
}

func (b *Repository) ForEach(bucketNames []string, f func(value []byte) error) error {
	return b.LoBucketFunc(bucketNames, func(bucket *bolt.Bucket) error {
		return bucket.ForEach(func(k, v []byte) error {
			return f(v)
		})
	})
}

// UpdateByID updates data by ID.
// if data which have ID does not exist, return error.
func (b *Repository) UpdateByID(bucketNames []string, data BoltData) error {
	return b.BucketFunc(bucketNames, func(bucket *bolt.Bucket) error {
		id := data.GetID()
		if id == 0 { // FIXME: implement HasID() to BoltData
			return fmt.Errorf("failed to update data of bolt. ID does not provided")
		}

		if bucket.Get(Itob(id)) == nil {
			return fmt.Errorf("failed to update data of bolt. provided ID(%d) does not exist", id)
		}

		return putByID(bucket, data)
	})
}

// SaveByID add or update data.
// If data ID does not provided, add new data.
// If data ID provided, update data which have the ID.
// If data ID provided, but data which have the ID does not exist, return error.
func (b *Repository) SaveByID(bucketNames []string, data BoltData) (retID uint64, err error) {
	err = b.BucketFunc(bucketNames, func(bucket *bolt.Bucket) error {
		retID = data.GetID()

		// update data
		if retID != 0 {
			if bucket.Get(Itob(retID)) == nil {
				return fmt.Errorf("failed to save data by ID(%d). ID is provided but does not exist", retID)
			}
			return putByID(bucket, data)
		}

		// add data
		id, err := addWithID(bucket, data)
		retID = id
		return err
	})
	return
}

func (b *Repository) Delete(bucketNames []string, id uint64) error {
	return b.BucketFunc(bucketNames, func(bucket *bolt.Bucket) error {
		return bucket.Delete(Itob(id))
	})
}

func (b *Repository) Close() error {
	return b.bolt.Close()
}

// ---- batch operations ----

func (b *Repository) BatchGetByString(bucketNames []string, keys []string) (dataList [][]byte, err error) {
	err = b.LoBucketFunc(bucketNames, func(bucket *bolt.Bucket) error {
		for _, key := range keys {
			dataList = append(dataList, bucket.Get([]byte(key)))
		}
		return nil
	})
	return dataList, err
}

// BatchUpdateByID update data by ID. If ID does not exist in bucket, skip the data.
func (b *Repository) BatchUpdateByID(bucketNames []string, dataList []BoltData) (updatedDataList []BoltData, skippedDataList []BoltData, err error) {
	err = b.BucketFunc(bucketNames, func(bucket *bolt.Bucket) error {
		for _, data := range dataList {
			key := Itob(data.GetID())
			if v := bucket.Get(key); v == nil {
				skippedDataList = append(skippedDataList, data)
				continue
			}

			s, err := json.Marshal(data)
			if err != nil {
				return fmt.Errorf("failed to marshal tag to json: %w", err)
			}

			if err := bucket.Put(key, s); err != nil {
				return err
			}
			updatedDataList = append(updatedDataList, data)
		}
		return nil
	})
	return
}

func (b *Repository) BatchGet(bucketNames []string, idList []uint64) (dataList [][]byte, err error) {
	err = b.LoBucketFunc(bucketNames, func(bucket *bolt.Bucket) error {
		for _, id := range idList {
			dataList = append(dataList, bucket.Get(Itob(id)))
		}
		return nil
	})
	return dataList, err
}

func (b *Repository) BatchAddWithID(bucketNames []string, dataList []BoltData) (idList []uint64, err error) {
	e := b.batchBucketFunc(bucketNames, func(bucket *bolt.Bucket) error {
		for _, data := range dataList {
			id, err := addWithID(bucket, data)
			if err != nil {
				return err
			}
			idList = append(idList, id)
		}
		return nil
	})
	return idList, e
}

func (b *Repository) BatchAddIntByString(bucketNames []string, keys []string, values []uint64) error {
	return b.batchBucketFunc(bucketNames, func(bucket *bolt.Bucket) error {
		for i, key := range keys {
			if err := bucket.Put([]byte(key), Itob(values[i])); err != nil {
				return err
			}
		}
		return nil
	})
}
