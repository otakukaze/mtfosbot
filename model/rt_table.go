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

// AddRT - add facebook line rt
func (p *LineFacebookRT) AddRT() (err error) {
	_, err = x.NamedExec(`insert into "public"."line_fb_rt" ("line", "facebook", "tmpl") values (:line, :facebook, :tmpl)`, p)
	return
}

// AddRT - add twitch line rt
func (p *LineTwitchRT) AddRT() (err error) {
	_, err = x.NamedExec(`insert into "public"."line_twitch_rt" ("line", "twitch", "type", "tmpl") values (:line, :twitch, :type, :tmpl)`, p)
	return
}

// AddRT - add youtube line rt
func (p *LineYoutubeRT) AddRT() (err error) {
	_, err = x.NamedExec(`insert into "public"."line_youtube_rt" ("line", "youtube", "tmpl") values (:line, :youtube, :tmpl)`, p)
	return
}
