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

// CheckGroup -
func CheckGroup(g string) (exists bool, err error) {
	ss := struct {
		C int `db:"c"`
	}{}

	err = x.Get(&ss, `select count(*) as c from "public"."line_group" where "id" = $1`, g)
	if err != nil {
		return false, err
	}
	return ss.C > 0, nil
}
