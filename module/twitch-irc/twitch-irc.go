package twitchirc

import (
	"fmt"
	"net"

	"gopkg.in/irc.v2"

	"git.trj.tw/golang/mtfosbot/module/config"
)

var client *irc.Client

// InitIRC -
func InitIRC() (err error) {
	conf := config.GetConf()
	conn, err := net.Dial("tcp", conf.Twitch.ChatHost)
	if err != nil {
		return
	}
	config := irc.ClientConfig{
		Handler: irc.HandlerFunc(ircHandle),
	}

	client = irc.NewClient(conn, config)

	err = client.Run()
	return
}

func ircHandle(c *irc.Client, m *irc.Message) {
	fmt.Println(m.String())
}
