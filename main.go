package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"

	"github.com/dannyhinshaw/big-dater/dater"
	"github.com/dannyhinshaw/big-dater/db"
)

func main() {
	var err error
	var id int64

	var origin dater.Origin
	var target string

	now := time.Now()

	dbEST, err := sql.Open("postgres", db.ConnEST)
	if err != nil {
		log.Fatalln("unable to connect to est database:", err)
	}

	dbPST, err := sql.Open("postgres", db.ConnPST)
	if err != nil {
		log.Fatalln("unable to connect to pst database:", err)
	}

	dbUTC, err := sql.Open("postgres", db.ConnUTC)
	if err != nil {
		log.Fatalln("unable to connect to utc database:", err)
	}

	storeEST := dater.NewStore(dbEST)
	if err := storeEST.CreateTables(); err != nil {
		log.Fatalln("error creating tables for est db:", err)
	}

	storePST := dater.NewStore(dbPST)
	if err := storePST.CreateTables(); err != nil {
		log.Fatalln("error creating tables for pst db:", err)
	}

	storeUTC := dater.NewStore(dbUTC)
	if err := storeUTC.CreateTables(); err != nil {
		log.Fatalln("error creating tables for utc db:", err)
	}

	/*
		EST Data
	*/

	origin = dater.OriginEST
	target = db.TargetStr(origin, dater.WithZone{})

	id, err = storeEST.InsertWithZone(dater.OriginEST, now)
	if err != nil {
		log.Fatalf("error inserting %s row: %s", target, err)
	}

	estWithZone, err := storeEST.GetWithZoneByID(id)
	if err != nil {
		log.Fatalf("error getting %s row with id %d: %s", target, id, err)
	}

	target = db.TargetStr(origin, dater.NoZone{})

	id, err = storeEST.InsertNoZone(dater.OriginEST, now)
	if err != nil {
		log.Fatalf("error inserting %s row: %s", target, err)
	}

	estNoZone, err := storeEST.GetNoZoneByID(id)
	if err != nil {
		log.Fatalf("error getting %s row with id %d: %s", target, id, err)
	}

	log.Printf("estWithZone: %+v", estWithZone)
	log.Printf("estNoZone::: %+v", estNoZone)

	/*
		PST Data
	*/

	origin = dater.OriginPST
	target = db.TargetStr(origin, dater.WithZone{})

	id, err = storePST.InsertWithZone(origin, now)
	if err != nil {
		log.Fatalf("error inserting %s row: %s", target, err)
	}

	pstWithZone, err := storePST.GetWithZoneByID(id)
	if err != nil {
		log.Fatalf("error getting %s row with id %d: %s", target, id, err)
	}

	target = db.TargetStr(origin, dater.NoZone{})

	id, err = storePST.InsertNoZone(origin, now)
	if err != nil {
		log.Fatalf("error inserting %s row: %s", target, err)
	}

	pstNoZone, err := storePST.GetNoZoneByID(id)
	if err != nil {
		log.Fatalf("error getting %s row with id %d: %s", target, id, err)
	}

	log.Printf("pstWithZone: %+v", pstWithZone)
	log.Printf("pstNoZone::: %+v", pstNoZone)

	/*
		UTC Data
	*/

	origin = dater.OriginUTC
	target = db.TargetStr(origin, dater.WithZone{})

	id, err = storeUTC.InsertWithZone(origin, now)
	if err != nil {
		log.Fatalf("error inserting %s row: %s", target, err)
	}

	utcWithZone, err := storeUTC.GetWithZoneByID(id)
	if err != nil {
		log.Fatalf("error getting %s row with id %d: %s", target, id, err)
	}

	target = db.TargetStr(origin, dater.NoZone{})

	id, err = storeUTC.InsertNoZone(origin, now)
	if err != nil {
		log.Fatalf("error inserting %s row: %s", target, err)
	}

	utcNoZone, err := storeUTC.GetNoZoneByID(id)
	if err != nil {
		log.Fatalf("error getting %s row with id %d: %s", target, id, err)
	}

	log.Printf("utcWithZone: %+v", utcWithZone)
	log.Printf("utcNoZone::: %+v", utcNoZone)
}
