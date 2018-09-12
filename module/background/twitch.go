package background

import (
	"fmt"
	"strings"

	"git.trj.tw/golang/mtfosbot/model"
	"git.trj.tw/golang/mtfosbot/module/apis/line"
	"git.trj.tw/golang/mtfosbot/module/apis/twitch"
)

func getStreamStatus() {
	fmt.Println("run twitch check")
	channels, err := model.GetAllTwitchChannel()
	if err != nil {
		return
	}
	var ids []string
	for _, v := range channels {
		ids = append(ids, v.ID)
	}

	info := twitch.GetUserStreamStatus(ids)
	fmt.Printf("info len: %d\n", len(info))
	if len(info) == 0 {
		return
	}
	for _, v := range info {
		for _, ch := range channels {
			if v.UserID == ch.ID {
				go checkStream(ch, v)
			}
		}
	}
}

func checkStream(ch *model.TwitchChannel, info *twitch.StreamInfo) {
	if info.ID == ch.LastStream {
		return
	}
	err := ch.GetGroups()
	if err != nil {
		return
	}
	err = ch.UpdateStream(info.ID)
	if err != nil {
		return
	}

	link := fmt.Sprintf("https://twitch.tv/%s", ch.Name)
	for _, v := range ch.Groups {
		if v.Notify {
			tmpl := v.Tmpl
			if len(tmpl) > 0 {
				tmpl = strings.Replace(tmpl, "{txt}", info.Title, -1)
				tmpl = strings.Replace(tmpl, "{link}", link, -1)
			} else {
				tmpl = fmt.Sprintf("%s\n%s", info.Title, link)
			}
			msg := line.TextMessage{
				Text: tmpl,
			}
			line.PushMessage(v.ID, msg)
		}
	}
}
