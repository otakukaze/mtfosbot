package model

import (
	"fmt"
	"log"

	"git.trj.tw/golang/mtfosbot/module/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var x *sqlx.DB

// NewDB - connect to database
func NewDB() (*sqlx.DB, error) {
	var err error
	conf := config.GetConf()
	connStr := fmt.Sprintf("user=%s password=%s host=%s sslmode=disable dbname=%s port=%d", conf.Database.User, conf.Database.Pass, conf.Database.Host, conf.Database.DB, conf.Database.Port)
	x, err = sqlx.Connect("postgres", connStr)
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

// GetDB -
func GetDB() *sqlx.DB {
	return x
}
