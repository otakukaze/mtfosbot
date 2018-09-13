package model

import (
	"database/sql"
	"errors"
	"time"
)

// KeyCommands - struct
type KeyCommands struct {
	Key     string    `db:"key" cc:"key"`
	Group   string    `db:"group" cc:"group"`
	Message string    `db:"message" cc:"message"`
	Ctime   time.Time `db:"ctime" cc:"ctime"`
	Mtime   time.Time `db:"mtime" cc:"ctime"`
}

// GetKeyCommand -
func GetKeyCommand(c, g string) (cmd *KeyCommands, err error) {
	if len(c) == 0 {
		return nil, errors.New("command is empty")
	}
	tmpCmd := struct {
		KeyCommands
		Message2 sql.NullString `db:"message2"`
	}{}
	query := `select  c.*, c2.message as message2 from "public"."key_commands" c
		left join "public"."key_commands" c2
		on c2.key = c.key and c2."group" = $2
		where c."key" = $1
		and c."group" = ''`
	err = x.Get(&tmpCmd, query, c, g)
	if err != nil {
		return nil, err
	}

	cmd = &tmpCmd.KeyCommands
	if tmpCmd.Message2.Valid {
		cmd.Message = tmpCmd.Message2.String
	}

	return
}
