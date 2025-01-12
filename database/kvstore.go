package database

import (
	"fmt"

	"go.etcd.io/bbolt"
)

type KVStore struct {
	db *bbolt.DB
}

// OpenKVStore opens or creates a key-value store database file
func OpenKVStore(filename string) (*KVStore, error) {
	db, err := bbolt.Open(filename, 0600, nil)
	if err != nil {
		return nil, err
	}
	return &KVStore{db: db}, nil
}

// CreateBucketIfNotExists creates a bucket if it does not exist
func (kv *KVStore) CreateBucketIfNotExists(bucketName string) error {
	fn := func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		return err
	}
	return kv.db.Update(fn)
}

// Put inserts a key-value pair into the bucket
func (kv *KVStore) Put(bucketName string, key, value []byte) error {
	fn := func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("bucket %s does not exist", bucketName)
		}
		return bucket.Put(key, value)
	}
	return kv.db.Update(fn)
}

// Get retrieves a value from the bucket by key
func (kv *KVStore) Get(bucketName string, key []byte) ([]byte, error) {
	var value []byte
	fn := func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("bucket %s does not exist", bucketName)
		}
		value = bucket.Get(key)
		return nil
	}
	err := kv.db.View(fn)
	return value, err
}

// Delete removes a key-value pair from the bucket
func (kv *KVStore) Delete(bucketName string, key []byte) error {
	fn := func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("bucket %s does not exist", bucketName)
		}
		return bucket.Delete(key)
	}
	return kv.db.Update(fn)
}

// CloseDatabase closes the database
func (kv *KVStore) Close() {
	kv.db.Close()
}
