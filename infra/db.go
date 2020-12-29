package infra

import (
	"fmt"

	bolt "go.etcd.io/bbolt"
)

func NewBoltDB(dbPath string) (*bolt.DB, error) {
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open db file from %s: %w", dbPath, err)
	}
	return db, nil
}
