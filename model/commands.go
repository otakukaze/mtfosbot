package model

import "time"

// Commands - struct
type Commands struct {
	Cmd     string    `db:"cmd" cc:"cmd"`
	Message string    `db:"message" cc:"message"`
	Group   string    `db:"group" cc:"group"`
	Ctime   time.Time `db:"ctime" cc:"ctime"`
	Mtime   time.Time `db:"mtime" cc:"ctime"`
}
