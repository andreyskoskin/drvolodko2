package datasource

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/andreyskoskin/drvolodko2/datamodel"
	"github.com/andreyskoskin/drvolodko2/datasource/auditor"
	"github.com/andreyskoskin/drvolodko2/datasource/auditprogram"
)

type (
	PostgresConfig struct {
		DBName   string `toml:"dbname"`
		Host     string `toml:"host"`
		Port     int    `toml:"port"`
		User     string `toml:"user"`
		Password string `toml:"password"`
	}

	Postgres struct {
		db *sql.DB
	}
)

func NewPostgres(config PostgresConfig) (_ *Postgres, err error) {
	var db *sql.DB
	if db, err = sql.Open("postgres", config.ConnectionString()); err != nil {
		return nil, err
	}
	return &Postgres{db: db}, nil
}

func (pg *Postgres) AuditPrograms() datamodel.AuditPrograms {
	return auditprogram.NewPostgres(pg.db)
}

func (pg *Postgres) Auditors() datamodel.Auditors {
	return auditor.NewPostgres(pg.db)
}

func (pg *Postgres) Close() error {
	return pg.db.Close()
}

func (c *PostgresConfig) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.DBName)
}
