package db

import (
	"github.com/dgraph-io/badger/v3"
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
