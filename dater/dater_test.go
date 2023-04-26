package dater_test

import (
	"database/sql"
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-txdb"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"

	"github.com/dannyhinshaw/big-dater/dater"
	"github.com/dannyhinshaw/big-dater/db"
)

const (
	txdbKeyEST = "txdbEST"
	txdbKeyPST = "txdbPST"
	txdbKeyUTC = "txdbUTC"
)

type PostgresTimezonesTestSuite struct {
	suite.Suite
	storeEST *dater.Store
	storePST *dater.Store
	storeUTC *dater.Store
}

func TestPostgresTimezonesTestSuite(t *testing.T) {
	txdb.Register(txdbKeyEST, "postgres", db.ConnEST)
	txdb.Register(txdbKeyPST, "postgres", db.ConnPST)
	txdb.Register(txdbKeyUTC, "postgres", db.ConnUTC)

	suite.Run(t, new(PostgresTimezonesTestSuite))
}

func (suite *PostgresTimezonesTestSuite) SetupTest() {
	dbEST, err := sql.Open(txdbKeyEST, db.ConnEST)
	if err != nil {
		log.Fatalln("unable to connect to est database:", err)
	}

	dbPST, err := sql.Open(txdbKeyPST, db.ConnPST)
	if err != nil {
		log.Fatalln("unable to connect to pst database:", err)
	}

	dbUTC, err := sql.Open(txdbKeyUTC, db.ConnUTC)
	if err != nil {
		log.Fatalln("unable to connect to utc database:", err)
	}

	suite.storeEST = dater.NewStore(dbEST)
	if err := suite.storeEST.CreateTables(); err != nil {
		log.Fatalln("error creating tables for est db:", err)
	}

	suite.storePST = dater.NewStore(dbPST)
	if err := suite.storePST.CreateTables(); err != nil {
		log.Fatalln("error creating tables for pst db:", err)
	}

	suite.storeUTC = dater.NewStore(dbUTC)
	if err := suite.storeUTC.CreateTables(); err != nil {
		log.Fatalln("error creating tables for utc db:", err)
	}
}

func (suite *PostgresTimezonesTestSuite) Test_Timezones() {
	var err error
	var id int64

	var origin dater.Origin
	var target string

	now := time.Now()

	/*
		EST Data
	*/

	origin = dater.OriginEST
	target = db.TargetStr(origin, dater.WithZone{})

	id, err = suite.storeEST.InsertWithZone(dater.OriginEST, now)
	if err != nil {
		log.Fatalf("error inserting %s row: %s", target, err)
	}

	estWithZone, err := suite.storeEST.GetWithZoneByID(id)
	if err != nil {
		log.Fatalf("error getting %s row with id %d: %s", target, id, err)
	}

	target = db.TargetStr(origin, dater.NoZone{})

	id, err = suite.storeEST.InsertNoZone(dater.OriginEST, now)
	if err != nil {
		log.Fatalf("error inserting %s row: %s", target, err)
	}

	estNoZone, err := suite.storeEST.GetNoZoneByID(id)
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

	id, err = suite.storePST.InsertWithZone(origin, now)
	if err != nil {
		log.Fatalf("error inserting %s row: %s", target, err)
	}

	pstWithZone, err := suite.storePST.GetWithZoneByID(id)
	if err != nil {
		log.Fatalf("error getting %s row with id %d: %s", target, id, err)
	}

	target = db.TargetStr(origin, dater.NoZone{})

	id, err = suite.storePST.InsertNoZone(origin, now)
	if err != nil {
		log.Fatalf("error inserting %s row: %s", target, err)
	}

	pstNoZone, err := suite.storePST.GetNoZoneByID(id)
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

	id, err = suite.storeUTC.InsertWithZone(origin, now)
	if err != nil {
		log.Fatalf("error inserting %s row: %s", target, err)
	}

	utcWithZone, err := suite.storeUTC.GetWithZoneByID(id)
	if err != nil {
		log.Fatalf("error getting %s row with id %d: %s", target, id, err)
	}

	target = db.TargetStr(origin, dater.NoZone{})

	id, err = suite.storeUTC.InsertNoZone(origin, now)
	if err != nil {
		log.Fatalf("error inserting %s row: %s", target, err)
	}

	utcNoZone, err := suite.storeUTC.GetNoZoneByID(id)
	if err != nil {
		log.Fatalf("error getting %s row with id %d: %s", target, id, err)
	}

	log.Printf("utcWithZone: %+v", utcWithZone)
	log.Printf("utcNoZone::: %+v", utcNoZone)
}
