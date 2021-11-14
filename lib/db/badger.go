package db

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/dustin/go-humanize"
	"github.com/maxkulish/dnscrypt-list/lib/logger"
	"go.uber.org/zap"
)

// Conn a wrapper for DB connection
type Conn struct {
	conn *badger.DB
}

// NewConn create new connection with the Badger database
// returns the pointer to the connection
func NewConn(path string) (*Conn, error) {

	db, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		return nil, err
	}

	return &Conn{
		conn: db,
	}, nil
}

// Close closes a DB
func (c *Conn) Close() {
	err := c.conn.Close()
	if err != nil {
		logger.Error("database closing error", zap.Error(err))
	}
}

// AddBunch takes map and adds them to the database
func (c *Conn) AddBunch(updates map[string]string) error {

	var err error
	txn := c.conn.NewTransaction(true)
	for k, v := range updates {
		if err := txn.Set([]byte(k), []byte(v)); err == badger.ErrTxnTooBig {
			err = txn.Commit()
			txn = c.conn.NewTransaction(true)
			_ = txn.Set([]byte(k), []byte(v))
		}
	}
	err = txn.Commit()

	logger.Debug(
		"bunch saved to the db",
		zap.String("domains", humanize.Comma(int64(len(updates)))),
	)
	return err
}

//GetAllKeys iterate all keys and returns them in the slice
func (c *Conn) GetAllKeys() []string {

	var keys []string
	err := c.conn.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()

			keys = append(keys, string(k))
		}
		return nil
	})

	if err != nil {
		logger.Error("get all keys error", zap.Error(err))
	}

	return keys
}
