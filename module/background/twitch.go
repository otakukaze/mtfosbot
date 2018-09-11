package background

import (
	"fmt"

	"git.trj.tw/golang/mtfosbot/module/utils"

	"git.trj.tw/golang/mtfosbot/model"
	"git.trj.tw/golang/mtfosbot/module/apis/twitch"
)

func getStreamStatus() {
	fmt.Println("run twitch check")
	channels, err := model.GetAllChannel()
	if err != nil {
		return
	}
	var ids []string
	for _, v := range channels {
		ids = append(ids, v.ID)
	}

	info := twitch.GetUserStreamStatus(ids)
	fmt.Printf("info len: %d\n", len(info))
	for _, v := range info {
		fmt.Println(utils.ToMap(v))
	}
}

func checkStream(ch *model.TwitchChannel) {

}
