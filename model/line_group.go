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

// CheckGroupOwner -
func CheckGroupOwner(user, g string) (exists bool, err error) {
	ss := struct {
		C int `db:"c"`
	}{}

	err = x.Get(&ss, `select count(*) as c from "public"."line_group" where "id" = $1 and "owner" = $2`, g, user)
	if err != nil {
		return false, err
	}
	return ss.C > 0, nil
}

// GetLineGroup -
func GetLineGroup(id string) (g *LineGroup, err error) {
	err = x.Get(&g, `select * from "public"."line_group" where "id" = $1`, id)
	return
}

// AddLineGroup -
func AddLineGroup(name, owner string, notify bool) (g *LineGroup, err error) {
	err = x.Get(&g, `insert into "public"."line_group" ("name", "owner", "notify") values ($1, $2, $3)`, name, owner, notify)
	if err != nil {
		return nil, err
	}
	return
}

// DeleteGroup -
func (p *LineGroup) DeleteGroup() (err error) {
	_, err = x.Exec(`delete from "public"."line_group" where "id" = $1`, p.ID)
	return
}
