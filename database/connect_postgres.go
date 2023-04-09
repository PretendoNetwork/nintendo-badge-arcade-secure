package database

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
	"github.com/PretendoNetwork/badge-arcade-secure/globals"
)

var postgres *sql.DB

func connectPostgres() {
	// Connect to Postgres

	var err error

	postgres, err = sql.Open("postgres", os.Getenv("DATABASE_URI"))
	if err != nil {
		globals.Logger.Critical(err.Error())
		return
	}

	globals.Logger.Success("Connected to Postgres!")

	_, err = postgres.Exec(`CREATE SCHEMA IF NOT EXISTS pretendo_badge_arcade`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return
	}

	globals.Logger.Success("Postgres schema created")

	// Create tables if missing

	_, err = postgres.Exec(`CREATE TABLE IF NOT EXISTS pretendo_badge_arcade.free_play_data (
			data_id int PRIMARY KEY,
			owner_id int,
			meta_binary bytea,
			created_time bigint,
			updated_time bigint,
			period smallint,
			flag int,
			referred_time bigint
		)`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return
	}

	_, err = postgres.Exec(`CREATE TABLE IF NOT EXISTS pretendo_badge_arcade.user_play_info (
		data_id int PRIMARY KEY,
		pid int,
		slot smallint,
		version int
	)`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return
	}

	globals.Logger.Success("Postgres tables created")
}
