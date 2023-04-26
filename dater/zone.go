package dater

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// zone is the same schema for both no_zone and with_zone tables
type zone struct {
	ID          int64
	DBOrigin    Origin
	DBTimezone  string
	DBTimestamp time.Time
	GoTimeNow   time.Time
	GoTimeUTC   time.Time
}

// NoZone is the schema for the `no_zone` tables.
type NoZone zone

func (z NoZone) TableName() string {
	return "no_zone"
}

// WithZone is the schema for the `with_zone` tables.
type WithZone zone

func (z WithZone) TableName() string {
	return "with_zone"
}

func (s Store) InsertNoZone(dbOrigin Origin, now time.Time) (int64, error) {
	var id int64
	err := s.db.QueryRow(sqlInsertNoZone, dbOrigin, now, now.UTC()).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error creating new no_zone zone: %w", err)
	}

	return id, nil
}

func (s Store) GetNoZoneByID(id int64) (*NoZone, error) {
	var r NoZone
	err := s.db.QueryRow(sqlSelectNoZoneByID, id).Scan(
		&r.ID, &r.DBTimezone,
		&r.DBOrigin, &r.DBTimestamp,
		&r.GoTimeNow, &r.GoTimeUTC,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error scanning %s row: %w", r.TableName(), err)
	}

	return &r, nil
}

func (s Store) InsertWithZone(dbOrigin Origin, now time.Time) (int64, error) {
	var id int64
	err := s.db.QueryRow(sqlInsertWithZone, dbOrigin, now, now.UTC()).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error creating new with_zone zone: %w", err)
	}

	return id, nil
}

func (s Store) GetWithZoneByID(id int64) (*WithZone, error) {
	var r WithZone
	err := s.db.QueryRow(sqlSelectWithZoneByID, id).Scan(
		&r.ID, &r.DBTimezone,
		&r.DBOrigin, &r.DBTimestamp,
		&r.GoTimeNow, &r.GoTimeUTC,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error scanning %s row: %w", r.TableName(), err)
	}

	return &r, nil
}

const (
	sqlInsertNoZone = `
		insert into no_zone (db_tz, db_origin, go_time_now, go_time_utc)
		values ((select current_setting('TIMEZONE')), $1, $2, $3)
		returning id`

	sqlSelectNoZoneByID = `
		select id, db_tz, db_origin, db_timestamp, go_time_now, go_time_utc 
		from no_zone
		where id = $1`

	sqlInsertWithZone = `
		insert into with_zone (db_tz, db_origin, go_time_now, go_time_utc)
		values ((select current_setting('TIMEZONE')), $1, $2, $3)
		returning id`

	sqlSelectWithZoneByID = `
		select id, db_tz, db_origin, db_timestamp, go_time_now, go_time_utc 
		from with_zone
		where id = $1`
)
