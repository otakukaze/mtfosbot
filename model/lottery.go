package model

import (
	"database/sql"
	"math/rand"
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
	offset := rand.Intn(10)
	err = x.Get(p, `select * from "public"."lottery" where "type" = $1 order by random() offset $2 limit 1`, t, offset)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return
}
