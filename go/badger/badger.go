package badger

import (
	"context"

	badger "github.com/dgraph-io/badger/v2"
)

// BadgerTokenManager is a TokenManager
// that writes and retrieves data from
// a badger kv store
type BadgerTokenManager struct {
	db *badger.DB
}

// NewBadgerTokenManager creates a new token manager persisting at the given
// path
func NewBadgerTokenManager(path string) (*BadgerTokenManager, error) {
	options := badger.DefaultOptions(path)
	options.Logger = nil
	db, err := badger.Open(options)
	if err != nil {
		return nil, err
	}
	return &BadgerTokenManager{db}, nil
}

// Set sets or updates the stored token
func (m *BadgerTokenManager) Set(ctx context.Context, subject, value string) error {
	return m.db.Update(func(tx *badger.Txn) error {
		return tx.Set([]byte(subject), []byte(value))
	})
}

// Get returns the stored token
func (m *BadgerTokenManager) Get(ctx context.Context, subject string) (string, error) {
	var value []byte
	if err := m.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(subject))
		if err != nil {
			return err
		}
		value, err = item.ValueCopy(nil)
		return err
	}); err != nil {
		return "", err
	}
	return string(value), nil
}

// Close closes the underlying badger database
func (m *BadgerTokenManager) Close() error {
	return m.db.Close()
}
