package model

import (
	"errors"
	"time"
)

// Commands - struct
type Commands struct {
	Cmd     string    `db:"cmd" cc:"cmd"`
	Message string    `db:"message" cc:"message"`
	Group   string    `db:"group" cc:"group"`
	Ctime   time.Time `db:"ctime" cc:"ctime"`
	Mtime   time.Time `db:"mtime" cc:"ctime"`
}

// GetAllCommands -
func GetAllCommands() (cmds []*Commands, err error) {
	err = x.Select(&cmds, `select * from "public"."commands"`)
	return
}

// GetGroupCommand -
func GetGroupCommand(c, g string) (cmd *Commands, err error) {
	if len(c) == 0 {
		return nil, errors.New("command is empty")
	}
	query := `select c.* from "public"."commands" c
		where c."cmd" = $1
		and (c."group" = '' or c."group" = $2)
		order by c."group" desc
		limit 1`
	row := x.QueryRowx(query, c, g)
	// err = x.Get(&cmd, query, c, g)
	cmd = &Commands{}
	err = row.StructScan(cmd)
	if err != nil {
		return nil, err
	}
	return
}
