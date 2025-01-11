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

func (kv *KVStore) CloseDatabase() {
	kv.db.Close()
}
