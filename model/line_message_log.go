package model

import "time"

// LineMessageLog -
type LineMessageLog struct {
	ID      string    `db:"id" cc:"id"`
	Group   string    `db:"group" cc:"group"`
	User    string    `db:"user" cc:"user"`
	Message string    `db:"message" cc:"message"`
	Ctime   time.Time `db:"ctime" cc:"ctime"`
	Mtime   time.Time `db:"mtime" cc:"mtime"`
}

// AddLineMessageLog -
func AddLineMessageLog(g, u, msg string) (msglog *LineMessageLog, err error) {
	query := `insert into "public"."line_message_log" ("group", "user", "message") values ($1, $2, $3)`
	err = x.Get(&msglog, query, g, u, msg)
	return
}
