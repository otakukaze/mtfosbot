package model

import "time"

// FacebookPage - struct
type FacebookPage struct {
	ID       string    `db:"id" cc:"id"`
	LastPost string    `db:"lastpost" cc:"lastpost"`
	Ctime    time.Time `db:"ctime" cc:"ctime"`
	Mtime    time.Time `db:"mtime" cc:"ctime"`
}
