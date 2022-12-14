package driver

import (
	"database/sql"
	"fmt"
	"time"
)

type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenDbConn = 25
const maxIdleDbConn = 25
const maxDbLifetime = 5 * time.Minute

func ConnectPostgres(dsn string) (*DB, error) {
	d, err := sql.Open("pgx", dsn)
	if err != nil {
		panic(err)
	}

	d.SetMaxOpenConns(maxOpenDbConn)
	d.SetMaxIdleConns(maxIdleDbConn)
	d.SetConnMaxLifetime(maxDbLifetime)

	dbConn.SQL = d

	if err := d.Ping(); err != nil {
		fmt.Println("Error!", err)
	}
	fmt.Println("*** pinged database successfully ***")

	return dbConn, err
}
