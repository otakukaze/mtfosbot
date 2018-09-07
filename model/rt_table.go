package model

type LineTwitchRT struct {
	Line   string `db:"line" cc:"line"`
	Twitch string `db:"twitch" cc:"twitch"`
	Tmpl   string `db:"tmpl" cc:"tmpl"`
	Type   string `db:"type" cc:"type"`
}

type LineFacebookRT struct {
	Line     string `db:"line" cc:"line"`
	Facebook string `db:"facebook" cc:"facebook"`
	Tmpl     string `db:"tmpl" cc:"tmpl"`
}

type LineYoutubeRT struct {
	Line    string `db:"line" cc:"line"`
	Youtube string `db:"youtube" cc:"youtube"`
	Tmpl    string `db:"tmpl" cc:"tmpl"`
}
