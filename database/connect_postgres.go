package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var postgres *sql.DB

func connectPostgres() {
	// Connect to Postgres

	var err error

	postgres, err = sql.Open("postgres", os.Getenv("DATABASE_URI"))
	if err != nil {
		panic(err)
	}

	_, err = postgres.Exec(`CREATE SCHEMA IF NOT EXISTS pretendo_badge_arcade`)
	if err != nil {
		fmt.Println("pretendo_badge_arcade")
		log.Fatal(err)
	}

	// Create tables if missing

	_, err = postgres.Exec(`CREATE TABLE IF NOT EXISTS pretendo_badge_arcade.free_play_data (
			data_id bigint PRIMARY KEY,
			owner_id int,
			meta_binary bytea,
			created_time bigint,
			updated_time bigint,
			period smallint,
			flag int,
			referred_time bigint
		)`)
	if err != nil {
		fmt.Println("pretendo_badge_arcade.free_play_data")
		log.Fatal(err)
	}

	_, err = postgres.Exec(`CREATE TABLE IF NOT EXISTS pretendo_badge_arcade.user_play_info (
		data_id bigint PRIMARY KEY,
		pid int,
		slot smallint,
		version int
	)`)
	if err != nil {
		fmt.Println("pretendo_badge_arcade.user_play_info")
		log.Fatal(err)
	}

	fmt.Println("Connected to Postgres")
}
