package persistence

import (
	"database/sql"
	"fmt"
	"time"
)

var Connections = map[string]*Connection{}

const (
	maxOpenDBConn = 10
	maxIdleConn   = 5
	maxDBLifeTime = 5 * time.Second
)

type Connection struct {
	provider provider
}

func NewConnection(connDetails *ConnectionDetails) (*Connection, error) {
	c := &Connection{}

	if newProvider, ok := Providers[connDetails.Provider]; ok {
		c.provider = newProvider(connDetails)
		return c, nil
	}

	return nil, fmt.Errorf("not valid provider %q", connDetails.Provider)
}

func (c *Connection) Open() (*sql.DB, error) {
	provider := c.provider
	db, err := sql.Open(string(provider.Driver()), provider.URL())
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenDBConn)
	db.SetMaxIdleConns(maxIdleConn)
	db.SetConnMaxLifetime(maxDBLifeTime)

	return db, nil
}

func (c *Connection) URL() string {
	return c.provider.URL()
}
