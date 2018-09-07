package model

import (
	"time"
)

// DonateSetting - struct
type DonateSetting struct {
	Twitch       string    `db:"twitch" cc:"twitch"`
	StartDate    time.Time `db:"start_date" cc:"start_date"`
	EndDate      time.Time `db:"end_date" cc:"end_date"`
	TargetAmount int       `db:"target_amount" cc:"target_amount"`
	Title        string    `db:"title" cc:"title"`
	Type         string    `db:"type" cc:"type"`
	Ctime        time.Time `db:"ctime" cc:"ctime"`
	Mtime        time.Time `db:"mtime" cc:"ctime"`
}
