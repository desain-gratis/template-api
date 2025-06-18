package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var POSTGRES_SUITE_API *sqlx.DB

// var CLICKHOUSE_SUITE_API

func GET_POSTGRES_SUITE_API() (*sqlx.DB, bool) {
	if POSTGRES_SUITE_API != nil {
		return POSTGRES_SUITE_API, true
	}

	var err error
	connString := fmt.Sprintf("user=%s dbname=%s sslmode=disable password=%s host=%s port=%s",
		CONFIG.GetString("postgres.db1.user"),
		CONFIG.GetString("postgres.db1.database_name"),
		CONFIG.GetString("postgres.db1.password"),
		CONFIG.GetString("postgres.db1.host"),
		CONFIG.GetString("postgres.db1.port"),
	)

	log.Debug().Msgf("Connecting to postgres database %v..", CONFIG.GetString("postgres.db1.database_name"))
	POSTGRES_SUITE_API, err = sqlx.Connect(
		"postgres",
		connString,
	)
	if err != nil {
		log.Fatal().Msgf("failed to connect postgres db: %v", err)
	}

	return POSTGRES_SUITE_API, POSTGRES_SUITE_API != nil
}
