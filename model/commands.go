package model

import (
	"errors"
	"fmt"
	"regexp"
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

// CommandsWithGroup -
type CommandsWithGroup struct {
	Commands
	GroupName string `db:"group_name" cc:"group_name"`
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

// GetCommandCount -
func GetCommandCount(where ...interface{}) (c int, err error) {
	query := `select count(*) as c from "public"."commands"`
	values := make([]interface{}, 0)
	if len(where) > 0 {
		if whereMap, ok := where[0].(map[string]string); ok {
			str := ""
			idx := 1
			for k, v := range whereMap {
				if len(str) > 0 {
					str += " and "
				}
				str += fmt.Sprintf(` "%s" = $%d `, k, idx)
				idx++
				values = append(values, v)
			}
			if len(str) > 0 {
				query += ` where ` + str
			}
		}
	}
	err = x.Get(&c, query, values)
	return
}

// GetCommands -
func GetCommands(where map[string]string, offset, limit int, order map[string]string) (cmds []*Commands, err error) {
	query := `select c.*, (case when g.name is null then '' else g.name end) as group_name from "public"."commands"
		left join "public"."line_group" g
		on g.id = c.group `
	values := make([]interface{}, (len(where) + len(order)))
	idx := 1

	if len(where) > 0 {
		str := ""
		for k, v := range where {
			if len(str) > 0 {
				str += " and "
			}
			str += fmt.Sprintf(` "%s" = $%d `, k, idx)
			idx++
			values = append(values, v)
		}
		if len(str) > 0 {
			query += ` where ` + str
		}
	}
	if offset >= 0 {
		query += fmt.Sprintf(" offset %d ", offset)
	}
	if limit > 0 {
		query += fmt.Sprintf(" limit %d ", limit)
	}

	if len(order) > 0 {
		regex, err := regexp.Compile("(?i)(desc|asc)")
		if err != nil {
			return nil, err
		}
		str := ""
		for k, v := range order {
			if !regex.Match([]byte(v)) {
				continue
			}
			if len(str) > 0 {
				str += " , "
			}
			str += fmt.Sprintf(` "%s" %s `, k, v)
		}
	}

	return
}
