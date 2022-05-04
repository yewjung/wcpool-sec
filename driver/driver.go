package driver

import (
	"database/sql"
	"log"
	"os"

	"github.com/lib/pq"
)

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var db *sql.DB

func ConnectDB() *sql.DB {
	pgurl, err := pq.ParseURL(os.Getenv("DB_URL"))
	logFatal(err)
	db, err = sql.Open("postgres", pgurl)
	logFatal(err)

	err = db.Ping()
	logFatal(err)

	return db
}
