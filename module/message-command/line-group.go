package msgcmd

import (
	"strconv"
	"strings"

	"git.trj.tw/golang/mtfosbot/module/apis/twitch"

	"git.trj.tw/golang/mtfosbot/model"
	lineobj "git.trj.tw/golang/mtfosbot/module/line-message/line-object"
)

func selectAct(cmd, sub, txt string, s *lineobj.SourceObject) (res string) {
	switch cmd {
	case "addgroup":
		return addLineGroup(sub, txt, s)
	case "addpage":
		return addFacebookPage(sub, txt, s)
	case "addtwitch":
		return addTwitchChannel(sub, txt, s)
	}
	return
}

func addLineGroup(sub, txt string, s *lineobj.SourceObject) (res string) {
	// args = groupName notify
	exists, err := model.CheckGroup(s.GroupID)
	if err != nil {
		return "run check group error"
	}
	if exists {
		return "group exists"
	}
	args := strings.Split(strings.Trim(txt, " "), " ")
	if len(args) < 2 {
		return "command args not match"
	}
	i, err := strconv.ParseInt(args[1], 10, 8)
	if err != nil || i < 0 || i > 1 {
		return "notify plases input 1 or 0"
	}

	_, err = model.AddLineGroup(args[0], s.UserID, i == 1)
	if err != nil {
		return "add group fail"
	}

	return "Success"
}

func addFacebookPage(sub, txt string, s *lineobj.SourceObject) (res string) {
	// args = pageid tmpl
	exists, err := model.CheckGroup(s.GroupID)
	if err != nil {
		return "run check group error"
	}
	if !exists {
		return "group not exists"
	}
	ok, err := model.CheckGroupOwner(s.UserID, s.GroupID)
	if err != nil {
		return "run check group owner fail"
	}
	if !ok {
		return "not owner"
	}

	args := strings.Split(strings.Trim(txt, " "), " ")
	if len(args) < 2 {
		return "command args not match"
	}

	page, err := model.GetFacebookPage(args[0])
	if err != nil {
		return "check facebook page fail"
	}
	if page == nil {
		page = &model.FacebookPage{
			ID: args[0],
		}
		err = page.AddPage()
		if err != nil {
			return "add facebook page fail"
		}
	}

	rt := &model.LineFacebookRT{
		Line:     s.GroupID,
		Facebook: args[0],
		Tmpl:     strings.Join(args[1:], " "),
	}

	err = rt.AddRT()
	if err != nil {
		return "add facebook page fail"
	}

	return "Success"
}

func addTwitchChannel(sub, txt string, s *lineobj.SourceObject) (res string) {
	// args = twitchLogin type tmpl
	exists, err := model.CheckGroup(s.GroupID)
	if err != nil {
		return "run check group error"
	}
	if !exists {
		return "group not exists"
	}
	ok, err := model.CheckGroupOwner(s.UserID, s.GroupID)
	if err != nil {
		return "run check group owner fail"
	}
	if !ok {
		return "not owner"
	}

	args := strings.Split(strings.Trim(txt, " "), " ")
	if len(args) < 3 {
		return "command args not match"
	}

	info := twitch.GetUserDataByName(args[0])
	if info == nil {
		return "get twitch user id fail"
	}

	ch := &model.TwitchChannel{
		ID:   info.ID,
		Name: info.DisplayName,
	}
	err = ch.Add()
	if err != nil {
		return "add twitch channel fail"
	}

	rt := &model.LineTwitchRT{
		Line:   s.GroupID,
		Twitch: info.ID,
		Tmpl:   strings.Join(args[2:], " "),
	}
	err = rt.AddRT()
	if err != nil {
		return "add rt data fail"
	}

	return "Success"
}
