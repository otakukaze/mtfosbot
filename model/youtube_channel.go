package model

import "time"

// YoutubeChannel -
type YoutubeChannel struct {
	ID        string    `db:"id" cc:"id"`
	Name      string    `db:"name" cc:"name"`
	LastVideo string    `db:"lastvideo" cc:"lastvideo"`
	Expire    int32     `db:"expire" cc:"expire"`
	Ctime     time.Time `db:"ctime" cc:"ctime"`
	Mtime     time.Time `db:"mtime" cc:"mtime"`
}
