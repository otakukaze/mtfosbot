package background

import (
	"git.trj.tw/golang/mtfosbot/model"
)

func checkOpay() {
	channels, err := model.GetAllTwitchChannel()
	if err != nil {
		return
	}
	for _, v := range channels {
		if len(v.OpayID) > 0 && v.Join {
			go getOpayData(v)
		}
	}
}

func getOpayData(ch *model.TwitchChannel) {

}
