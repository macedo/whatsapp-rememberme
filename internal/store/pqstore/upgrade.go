package pqstore

import (
	"database/sql"
	"log"
)

type upgradeFunc func(*sql.Tx, *Container) error

var Upgrades = [...]upgradeFunc{
	upgradeV1,
}

func (c *Container) Upgrade() error {
	version, err := c.getVersion()
	if err != nil {
		return err
	}

	for ; version < len(Upgrades); version++ {
		var tx *sql.Tx
		tx, err := c.DB.Begin()
		if err != nil {
			return err
		}

		migrateFunc := Upgrades[version]
		log.Printf("upgrading database to v%d", version+1)
		err = migrateFunc(tx, c)
		if err != nil {
			_ = tx.Rollback()
			return err
		}

		if err = c.setVersion(tx, version+1); err != nil {
			return err
		}

		if err = tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

func (c *Container) getVersion() (int, error) {
	_, err := c.DB.Exec("CREATE TABLE IF NOT EXISTS schema_migrations (version INTEGER)")
	if err != nil {
		return -1, err
	}

	version := 0
	row := c.DB.QueryRow("SELECT version FROM schema_migrations LIMIT 1")
	if row != nil {
		_ = row.Scan(&version)
	}

	return version, nil
}

func (c *Container) setVersion(tx *sql.Tx, version int) error {
	_, err := tx.Exec("DELETE FROM schema_migrations")
	if err != nil {
		return err
	}
	_, err = tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", version)
	return err
}

// - creates users table
func upgradeV1(tx *sql.Tx, c *Container) (err error) {
	_, err = tx.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`CREATE EXTENSION IF NOT EXISTS "citext"`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS users(
		id UUID DEFAULT uuid_generate_v4(),
		username CITEXT NOT NULL UNIQUE,
		encrypted_password VARCHAR NOT NULL,
		CONSTRAINT users_pkey PRIMARY KEY (id)
	)`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`INSERT INTO users (username, encrypted_password)
		VALUES ($1, $2)
	`, "admin", c.App.Encryptor.MustDigest("admin"))
	if err != nil {
		return err
	}

	return nil
}
