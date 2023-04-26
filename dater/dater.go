package dater

import (
	"database/sql"
	"fmt"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s Store) CreateTables() error {
	_, err := s.db.Exec(sqlCreateTableWithZone)
	if err != nil {
		return fmt.Errorf("error creating with_zone table: %w", err)
	}

	_, err = s.db.Exec(sqlCreateTableNoZone)
	if err != nil {
		return fmt.Errorf("error creating no_zone table: %w", err)
	}

	return nil
}

const (
	sqlCreateTableWithZone = `
		create table if not exists with_zone (
			id bigserial primary key,
			db_tz varchar,
			db_origin varchar,
			db_timestamp timestamptz default now(),
			go_time_now timestamptz,
			go_time_utc timestamptz
		)`

	sqlCreateTableNoZone = `
		create table if not exists no_zone (
			id bigserial primary key,
			db_tz varchar,
			db_origin varchar,
			db_timestamp timestamp default now(),
			go_time_now timestamp,
			go_time_utc timestamp
		)`
)
