package dbutils

import (
	"fmt"

	"github.com/dgraph-io/badger/v4"
)

type DB struct {
	db *badger.DB
}

func NewDB() (*DB, error) {
	opts := badger.DefaultOptions("/tmp/badger").WithLoggingLevel(badger.WARNING)
	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("> badger.Open: %v", err.Error())
	}
	return &DB{db: db}, nil
}

func (d *DB) Update(key, value []byte) error {
	return d.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

func (d *DB) Get(key []byte) ([]byte, error) {
	var valCopy []byte
	err := d.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		valCopy, err = item.ValueCopy(nil)
		return err
	})
	if err != nil {
		return nil, err
	}
	return valCopy, nil
}

func (d *DB) Delete(key []byte) error {
	return d.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
}

func (d *DB) Close() error {
	return d.db.Close()
}
