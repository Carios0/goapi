package data

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"log"
	"time"
)

func ConnectDB() *sql.DB {

	var err error
	// define database DSN
	cfg := mysql.Config{
		User:      "app",
		Passwd:    "app",
		Net:       "tcp",
		Addr:      "mysql:3306",
		DBName:    "app",
		ParseTime: true} // necessary for parsing SQL DATETIME into go time.Time

	// retry connection (else: using "docker-compose restart" golang container starts before mysql is ready)
	for i := 0; i < 10; i++ {

		log.Println("Starting...")
		db, err := sql.Open("mysql", cfg.FormatDSN())
		if err == nil {
			// ping db to test reachability
			err = db.Ping()
			if err == nil {
				log.Println("Db connected!")
				return db
			}
		}
		log.Printf("Failed to connect to database: %v. Retrying...", err)
		time.Sleep(1 * time.Second) // Wait for 1 seconds before retrying

	}
	log.Fatal("Could not connect to database:", err)
	return nil
}

// calling defer db.Close() directly in main closes it immediately and here it works fine, why?
func CloseDB(db *sql.DB) {
	defer db.Close()
}
