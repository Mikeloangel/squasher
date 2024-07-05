package config

import (
	"database/sql"

	"github.com/Mikeloangel/squasher/internal/logger"
)

// GetDB opens connection to db
func GetDB(conf *Config) *sql.DB {
	db, err := sql.Open("pgx", conf.DBDSN)
	if err != nil {
		logger.Fatal(err)
		return &sql.DB{}
	}

	return db
}
