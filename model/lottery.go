package model

import (
	"database/sql"
	"time"
)

// Lottery -
type Lottery struct {
	ID      string    `db:"id" cc:"id"`
	Type    string    `db:"type" cc:"type"`
	Message string    `db:"message" cc:"message"`
	Ctime   time.Time `db:"ctime" cc:"ctime"`
	Mtime   time.Time `db:"mtime" cc:"mtime"`
}

// GetRandomLotteryByType -
func GetRandomLotteryByType(t string) (p *Lottery, err error) {
	p = &Lottery{}
	err = x.Get(p, `select * from "public"."lottery" where "type" = $1 order random() limit 1`, t)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return
}
