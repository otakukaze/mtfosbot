package model

import (
	"database/sql"
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
	tmpCmd := struct {
		Commands
		Message2 sql.NullString `db:"message2"`
	}{}
	query := `select  c.*, c2.message as message2 from "public"."commands" c
		left join "public"."commands" c2
		on c2.cmd = c.cmd and c2."group" = $2
		where c."cmd" = $1
		and c."group" = ''`
	err = x.Get(&tmpCmd, query, c, g)
	if err != nil {
		return nil, err
	}

	cmd = &tmpCmd.Commands
	if tmpCmd.Message2.Valid {
		cmd.Message = tmpCmd.Message2.String
	}

	return
}
