package model

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

var x *sqlx.DB

// NewDB - connect to database
func NewDB() (*sqlx.DB, error) {
	var err error
	connStr := fmt.Sprintf("user=%s host=%s sslmode=disable dbname=%s port=%d", "postgres", "localhost", "mtfosbot", 5432)
	x, err := sqlx.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	x.SetMaxIdleConns(10)
	x.SetMaxOpenConns(200)
	err = x.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return x, err
}
