package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/macedo/whatsapp-rememberme/internal/infrastructure/persistence"
	"github.com/macedo/whatsapp-rememberme/pkg/env"
)

var ENV = env.Get("APP_ENV", "development")

var (
	down           *bool
	migrateCmd     = flag.NewFlagSet("migrate", flag.ExitOnError)
	migrationsPath = "file://internal/infrastructure/persistence/postgres/migrations"
)

var (
	ErrInvalidArgument = errors.New("invalid argument, expected command `migrate`")
)

func init() {
	down = migrateCmd.Bool("down", false, "run migrations down")
}

func main() {
	if err := run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) < 2 {
		return ErrInvalidArgument
	}

	if err := persistence.LoadConfigFile(); err != nil {
		return err
	}

	switch args[1] {
	case "migrate":
		migrateCmd.Parse(args[2:])
		conn := persistence.Connections[ENV]

		if *down {
			return migrateDown(conn)
		}

		return migrateUp(conn)

	default:
		return ErrInvalidArgument
	}
}

func migrateDown(conn *persistence.Connection) error {
	m, err := migrate.New(migrationsPath, conn.URL())
	if err != nil {
		return fmt.Errorf("error loading migrations - %w", err)
	}

	lastVersion, _, _ := m.Version()

	if err := m.Steps(-1); err != nil {
		m.Force(int(lastVersion))
		return fmt.Errorf("error migrate down - %w", err)
	}

	return nil
}

func migrateUp(conn *persistence.Connection) error {
	m, err := migrate.New(migrationsPath, conn.URL())
	if err != nil {
		return fmt.Errorf("error loading migrations - %w", err)
	}

	lastVersion, _, _ := m.Version()

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		m.Force(int(lastVersion))
		return fmt.Errorf("error migrate up - %w", err)
	}

	return nil
}
