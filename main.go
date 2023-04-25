package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	const (
		utcDBPort  = 7654
		eastDBPort = 6543
	)

	fmtConnStr := "port=%d user=postgres sslmode=disable"

	connUTC := fmt.Sprintf(fmtConnStr, utcDBPort)

	dbUTC, err := sql.Open("postgres", connUTC)
	if err != nil {
		log.Fatalln("error creating utc db connection:", err)
	}

	connEast := fmt.Sprintf(fmtConnStr, eastDBPort)

	dbEast, err := sql.Open("postgres", connEast)
	if err != nil {
		log.Fatalln("error creating east db connection:", err)
	}

	if err := createTables(dbUTC); err != nil {
		log.Fatalln("error creating tables for db utc:", err)
	}

	if err := createTables(dbEast); err != nil {
		log.Fatalln("error creating tables for db east:", err)
	}

	if err := insertLocalRows(dbUTC); err != nil {
		log.Fatalln("error inserting local rows in db utc:", err)
	}

	if err := insertLocalRows(dbEast); err != nil {
		log.Fatalln("error inserting local rows in db east:", err)
	}

	log.Println("big dater: success")
}

// insertLocalRows inserts rows with default values local to each postgres database.
func insertLocalRows(db *sql.DB) error {
	insertNoZone, err := db.Prepare(`insert into no_zones (go_time_now, go_time_utc) values ($1, $2) returning id`)
	if err != nil {
		return fmt.Errorf("error preparing no_zones insert statement: %w", err)
	}

	now := time.Now()
	if _, err := insertNoZone.Exec(now, now.UTC()); err != nil {
		return fmt.Errorf("error creating new no_zones row: %w", err)
	}

	insertWithZone, err := db.Prepare(`insert into with_zones (go_time_now, go_time_utc) values ($1, $2) returning id`)
	if err != nil {
		return fmt.Errorf("error preparing with_zones insert statement: %w", err)
	}

	now = time.Now()
	if _, err := insertWithZone.Exec(now, now.UTC()); err != nil {
		return fmt.Errorf("error creating new with_zones row: %w", err)
	}

	return nil
}

func createTables(db *sql.DB) error {

	_, err := db.Exec(`
	create table if not exists no_zones (
    	id bigserial primary key,
    	db_timestamp timestamp default now(),
    	go_time_now timestamp,
    	go_time_utc timestamp
	)`)
	if err != nil {
		return fmt.Errorf("error creating no_zones table: %w", err)
	}

	_, err = db.Exec(`
	create table if not exists with_zones (
    	id bigserial primary key,
    	db_timestamptz timestamptz default now(),
		go_time_now timestamptz,
    	go_time_utc timestamptz
	)`)
	if err != nil {
		return fmt.Errorf("error creating with_zones table: %w", err)
	}

	return nil
}
