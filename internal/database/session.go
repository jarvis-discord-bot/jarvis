package database

import (
	"database/sql"
	"fmt"
	"time"

	mysqldriver "github.com/go-sql-driver/mysql"
	"github.com/pedromsmoreira/jarvis/internal/configuration"
)

type SqlSession struct {
	Db          *sql.DB
	SqlSettings *configuration.Sql
}

func (recv *SqlSession) Close() error {
	return recv.Db.Close()
}

func Connect(settings *configuration.Sql) (*SqlSession, error) {
	mysqlCfg, err := mysqldriver.ParseDSN(settings.Dsn)
	mysqlCfg.InterpolateParams = true // recommended by planetscale
	mysqlCfg.MultiStatements = true   // required for golang-migrate
	if err != nil {
		return nil, fmt.Errorf("error initializing sql DSN: %w", err)
	}

	db, err := sql.Open("mysql", mysqlCfg.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("error initializing sql connection: %w", err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging sql connection: %w", err)
	}

	return &SqlSession{Db: db, SqlSettings: settings}, nil
}
