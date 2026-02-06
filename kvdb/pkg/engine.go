package pkg

import (
	"errors"
	"sync"
)

type KVDB struct {
	mu          sync.RWMutex
	collections Collections
}

func NewKVDB() *KVDB {
	return &KVDB{
		collections: make(Collections),
	}
}

var ErrCollectionAlreadyExists = errors.New("collection already exists")

func (db *KVDB) CreateCollection(collectionName string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, ok := db.collections[collectionName]; ok {
		return ErrCollectionAlreadyExists
	}

	db.collections[collectionName] = NewCollection()

	return nil
}

func (db *KVDB) GetCollection(collectionName string) (*Collection, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	collection, ok := db.collections[collectionName]
	if !ok {
		return nil, ErrCollectionNotFound
	}
	return collection, nil
}
