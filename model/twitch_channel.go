package model

import "time"

// TwitchGroup -
type TwitchGroup struct {
	*LineGroup
	Tmpl string `db:"tmpl"`
}

// TwitchChannel - struct
type TwitchChannel struct {
	ID         string         `db:"id" cc:"id"`
	Name       string         `db:"name" cc:"name"`
	LastStream string         `db:"laststream" cc:"laststream"`
	Join       bool           `db:"join" cc:"join"`
	OpayID     string         `db:"opayid" cc:"opayid"`
	Ctime      time.Time      `db:"ctime" cc:"ctime"`
	Mtime      time.Time      `db:"mtime" cc:"ctime"`
	Groups     []*TwitchGroup `db:"-"`
}

// GetAllTwitchChannel -
func GetAllTwitchChannel() (channels []*TwitchChannel, err error) {
	err = x.Select(&channels, `select * from "public"."twitch_channel"`)
	return
}

// GetJoinChatChannel -
func GetJoinChatChannel() (channels []*TwitchChannel, err error) {
	err = x.Select(&channels, `select * from "public"."twitch_channel" where "join" = true`)
	return
}

// Add -
func (p *TwitchChannel) Add() (err error) {
	rows, err := x.NamedQuery(`insert into "public"."twitch_channel" ("name", "laststream", "join", "opayid") values (:name, :laststream, :join, :opayid) returning *`, p)
	if err != nil {
		return err
	}
	err = rows.StructScan(p)
	return
}

// UpdateStream -
func (p *TwitchChannel) UpdateStream(streamID string) (err error) {
	query := `update "public"."twitch_channel" set "laststream" = $1 where "id" = $2`
	_, err = x.Exec(query, streamID, p.ID)
	if err != nil {
		return
	}
	p.LastStream = streamID
	return
}

// GetGroups -
func (p *TwitchChannel) GetGroups() (err error) {
	query := `select g.*, rt.tmpl as tmpl from "public"."twitch_channel" tw
		left join "public"."line_twitch_rt" rt
		on rt.twitch = tw.id
		left join "public"."line_group" g
		on g.id = rt.line
		where 
		tw.id = $1`
	err = x.Select(&p.Groups, query, p.ID)
	return
}
