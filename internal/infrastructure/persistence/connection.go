package persistence

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

var Connections = map[string]*Connection{}

var db *sql.DB

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

func DB() *sql.DB {
	if db == nil {
		log.Fatal(fmt.Errorf("database not initialized"))
	}

	return db
}

func (c *Connection) Open() error {
	var err error

	provider := c.provider
	db, err = sql.Open(string(provider.Driver()), provider.URL())
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(maxOpenDBConn)
	db.SetMaxIdleConns(maxIdleConn)
	db.SetConnMaxLifetime(maxDBLifeTime)

	return nil
}

func (c *Connection) URL() string {
	return c.provider.URL()
}
