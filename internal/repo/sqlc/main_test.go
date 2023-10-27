package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/msqtt/sevencowcloud-shortvideo-service/internal/pkg/config"
)

var(
	testQueries *Queries
	testDB *sql.DB
)

func TestMain(m *testing.M) {
	cnf, err := config.LoadConfig("../../../configs")
	if err != nil {
		log.Fatal(err)
	}
	testDB, err = sql.Open(cnf.DBDriver, cnf.DevDBSource)
	if err != nil {
		log.Fatal(err)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}
