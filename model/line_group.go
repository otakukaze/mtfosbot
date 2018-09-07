package model

import "time"

// LineGroup - struct
type LineGroup struct {
	ID     string    `db:"id" cc:"id"`
	Name   string    `db:"name" cc:"name"`
	Notify bool      `db:"notify" cc:"notify"`
	Owner  string    `db:"owner" cc:"owner"`
	Ctime  time.Time `db:"ctime" cc:"ctime"`
	Mtime  time.Time `db:"mtime" cc:"ctime"`
}
