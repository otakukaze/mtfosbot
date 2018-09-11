package model

import "time"

// TwitchChannel - struct
type TwitchChannel struct {
	ID         string    `db:"id" cc:"id"`
	Name       string    `db:"name" cc:"name"`
	LastStream string    `db:"laststream" cc:"laststream"`
	Join       bool      `db:"join" cc:"join"`
	OpayID     string    `db:"opayid" cc:"opayid"`
	Ctime      time.Time `db:"ctime" cc:"ctime"`
	Mtime      time.Time `db:"mtime" cc:"ctime"`
}

// GetAllChannel -
func GetAllChannel() (channels []*TwitchChannel, err error) {
	err = x.Select(&channels, `select * from "public"."twitch_channel"`)
	return
}
