package model

import "time"

// YoutubeGroup -
type YoutubeGroup struct {
	*LineGroup
	Tmpl string `db:"tmpl"`
}

// YoutubeChannel -
type YoutubeChannel struct {
	ID        string          `db:"id" cc:"id"`
	Name      string          `db:"name" cc:"name"`
	LastVideo string          `db:"lastvideo" cc:"lastvideo"`
	Expire    int64           `db:"expire" cc:"expire"`
	Ctime     time.Time       `db:"ctime" cc:"ctime"`
	Mtime     time.Time       `db:"mtime" cc:"mtime"`
	Groups    []*YoutubeGroup `db:"-"`
}

// GetYoutubeChannelWithID -
func GetYoutubeChannelWithID(id string) (yt *YoutubeChannel, err error) {
	err = x.Get(&yt, `select * from "public"."youtube_channel" where "id" = $1`, id)
	return
}

// Add -
func (p *YoutubeChannel) Add() (err error) {
	stmt, err := x.PrepareNamed(`insert into "public"."youtube_channel" ("id", "name") values (:id, :name) returning *`)
	if err != nil {
		return err
	}
	err = stmt.Get(p, p)
	return
}

// UpdateLastVideo -
func (p *YoutubeChannel) UpdateLastVideo(vid string) (err error) {
	p.LastVideo = vid
	_, err = x.NamedExec(`update "public"."youtube_channel" set "lastvideo" = :lastvideo where "id" = :id`, p)
	return
}

// UpdateExpire -
func (p *YoutubeChannel) UpdateExpire(t int64) (err error) {
	p.Expire = t
	_, err = x.NamedExec(`update "public"."youtube_channel" set "expire" = :expire where "id" = :id`, p)
	if err != nil {
		return err
	}
	return
}

// GetGroups -
func (p *YoutubeChannel) GetGroups() (err error) {
	query := `select g.*, rt.tmpl as tmpl from "public"."youtube_channel" yt 
		left join "public"."line_youtube_rt" rt 
		on rt.youtube = yt.id
		left join "public"."line_group" g
		on g.id = rt.line
		where 
		yt.id = $1`
	err = x.Select(&p.Groups, query, p.ID)
	return
}
