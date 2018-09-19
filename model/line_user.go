package model

import "time"

// LineUser -
type LineUser struct {
	ID    string    `db:"id" cc:"id"`
	Name  string    `db:"name" cc:"name"`
	Ctime time.Time `db:"ctime" cc:"ctime"`
	Mtime time.Time `db:"mtime" cc:"mtime"`
}

// GetLineUserByID -
func GetLineUserByID(id string) (u *LineUser, err error) {
	err = x.Get(&u, `select * from "public"."line_user" where "id" = $1`, id)
	return
}

// Add -
func (p *LineUser) Add() (err error) {
	_, err = x.NamedExec(`insert into "public"."line_user" ("id", "name") values (:id, :name)`, p)
	return
}
