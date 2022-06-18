package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"testing"
)

var testConnection *Queries
var dbTest *sql.DB

const (
	dbDriver = "mysql"
	dbPath   = "root:1234567@tcp(localhost:3308)/simple_bank?parseTime=true"
)

func TestMain(m *testing.M) {
	coon, err := sql.Open(dbDriver, dbPath)
	dbTest = coon
	if err != nil {
		log.Fatal("can not connect:", err)
	}
	testConnection = New(coon)

	os.Exit(m.Run())
}
