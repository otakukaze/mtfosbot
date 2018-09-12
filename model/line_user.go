package model

import "time"

// LineUser -
type LineUser struct {
	ID    string    `db:"id" cc:"id"`
	Name  string    `db:"name" cc:"name"`
	Ctime time.Time `db:"ctime" cc:"ctime"`
	Mtime time.Time `db:"mtime" cc:"mtime"`
}
