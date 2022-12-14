package sqlstore

import (
	"database/sql"
)

type upgradeFunc func(*sql.Tx, *Container) error

var Upgrades = [...]upgradeFunc{}

func (c *Container) Upgrade() error {
	version, err := c.getVersion()
	if err != nil {
		return err
	}

	for ; version < len(Upgrades); version++ {
		var tx *sql.Tx
		tx, err := c.db.Begin()
		if err != nil {
			return err
		}

		migrateFunc := Upgrades[version]
		c.log.Info().Int("version", version+1).Msg("upgrading database")
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
	_, err := c.db.Exec("CREATE TABLE IF NOT EXISTS schema_migrations (version INTEGER)")
	if err != nil {
		return -1, err
	}

	version := 0
	row := c.db.QueryRow("SELECT version FROM schema_migrations LIMIT 1")
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
	_, err = tx.Exec("INSERTT INTO schema_migrations (version) VALUES ($1)", version)
	return err
}
