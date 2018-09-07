package model

import "time"

// OpayDonateList - struct
type OpayDonateList struct {
	OpayID   string    `db:"opayid" cc:"opayid"`
	DonateID string    `db:"donate_id" cc:"donate_id"`
	Price    int       `db:"price" cc:"price"`
	Text     string    `db:"text" cc:"text"`
	Ctime    time.Time `db:"ctime" cc:"ctime"`
	Mtime    time.Time `db:"mtime" cc:"ctime"`
}
