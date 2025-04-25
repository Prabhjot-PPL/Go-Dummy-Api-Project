package persistance

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func ConnectToDatabase() (*Database, error) {

	db, err := sql.Open("postgres", "postgresql://postgres:password@localhost:5433/mydb?sslmode=disable")

	if err != nil {
		return nil, err
	}

	fmt.Printf("Connected to Database \n\n")

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &Database{db: db}, nil

}
