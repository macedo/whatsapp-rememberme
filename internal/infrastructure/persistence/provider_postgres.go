package persistence

import (
	"fmt"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const driverPGX = "pgx"

func init() {
	Providers["postgres"] = NewPostgresSQL
}

type postgres struct {
	ConnectionDetails *ConnectionDetails
}

func (p postgres) Driver() Driver {
	return driverPGX
}

func (p postgres) URL() string {
	connDetails := p.ConnectionDetails
	if url := connDetails.URL; url != "" {
		return url
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&timezone=UTC&connect_timeout=5",
		connDetails.User,
		connDetails.Password,
		connDetails.Host,
		connDetails.Port,
		connDetails.Database,
	)
}

func NewPostgresSQL(connDetails *ConnectionDetails) provider {
	return &postgres{
		ConnectionDetails: connDetails,
	}
}
