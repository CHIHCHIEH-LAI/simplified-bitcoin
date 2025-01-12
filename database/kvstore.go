package database

import "go.etcd.io/bbolt"

type KVStore struct {
	db *bbolt.DB
}

// OpenDatabase opens or creates a database file
func OpenDatabase(filename string) (*KVStore, error) {
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
		return bucket.Put(key, value)
	}
	return kv.db.Update(fn)
}

// Get retrieves a value from the bucket by key
func (kv *KVStore) Get(bucketName string, key []byte) ([]byte, error) {
	var value []byte
	fn := func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		value = bucket.Get(key)
		return nil
	}
	err := kv.db.View(fn)
	return value, err
}

// CloseDatabase closes the database
func (kv *KVStore) CloseDatabase() {
	kv.db.Close()
}
