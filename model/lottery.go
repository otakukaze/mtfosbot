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
	parr, err := GetRandomLotteryByTypeAndLimit(t, 1)
	if err != nil {
		return nil, err
	}
	if parr == nil {
		return nil, nil
	}
	return parr[0], err
}

// GetRandomLotteryByTypeAndLimit -
func GetRandomLotteryByTypeAndLimit(t string, limit int) (p []*Lottery, err error) {
	offset := rand.Intn(10)
	err = x.Select(&p, `select * from "public"."lottery" where "type" = $1 order by random() offset $2 limit $3`, t, offset, limit)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return
}
