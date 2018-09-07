package model

import "time"

// KeyCommands - struct
type KeyCommands struct {
	Key     string    `db:"key" cc:"key"`
	Group   string    `db:"group" cc:"group"`
	Message string    `db:"message" cc:"message"`
	Ctime   time.Time `db:"ctime" cc:"ctime"`
	Mtime   time.Time `db:"mtime" cc:"ctime"`
}
