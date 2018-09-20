package twitchirc

import (
	"crypto/tls"
	"fmt"
	"time"

	"git.trj.tw/golang/mtfosbot/model"

	"gopkg.in/irc.v2"

	"git.trj.tw/golang/mtfosbot/module/config"
)

var client *irc.Client
var queue *QueueList
var channels []string

// InitIRC -
func InitIRC() {
	conf := config.GetConf()
	tlsConf := &tls.Config{}
	conn, err := tls.Dial("tcp", conf.Twitch.ChatHost, tlsConf)
	// conn, err := net.Dial("tcp", conf.Twitch.ChatHost)
	if err != nil {
		return
	}
	defer conn.Close()

	channels = make([]string, 0)
	queue = NewQueue()
	runQueue()
	ReJoin()

	config := irc.ClientConfig{
		Nick:    conf.Twitch.BotUser,
		Pass:    conf.Twitch.BotOauth,
		Handler: irc.HandlerFunc(ircHandle),
	}

	client = irc.NewClient(conn, config)

	err = client.Run()
	if err != nil {
		fmt.Println("twitch chat connect fail")
	}
}

// SendMessage -
func SendMessage(ch, msg string) {
	if len(ch) == 0 {
		return
	}

	if indexOf(channels, ch) == -1 {
		return
	}

	m := &MsgObj{
		Command: "PRIVMSG",
		Params: []string{
			fmt.Sprintf("#%s", ch),
			fmt.Sprintf("%s", msg),
		},
	}
	queue.Add(m)
}

// ReJoin -
func ReJoin() {
	ch, err := model.GetAllTwitchChannel()
	if err != nil {
		return
	}
	LeaveAllChannel()
	for _, v := range ch {
		if v.Join {
			JoinChannel(v.Name)
		}
	}
}

// JoinChannel -
func JoinChannel(ch string) {
	if len(ch) == 0 {
		return
	}

	if indexOf(channels, ch) != -1 {
		return
	}

	m := &MsgObj{
		Command: "JOIN",
		Params: []string{
			fmt.Sprintf("#%s", ch),
		},
	}
	queue.Add(m)
}

// LeaveChannel -
func LeaveChannel(ch string) {
	if len(ch) == 0 {
		return
	}

	if indexOf(channels, ch) == -1 {
		return
	}

	m := &MsgObj{
		Command: "PART",
		Params: []string{
			fmt.Sprintf("#%s", ch),
		},
	}
	queue.Add(m)
}

// LeaveAllChannel -
func LeaveAllChannel() {
	if len(channels) == 0 {
		return
	}
	for _, v := range channels {
		m := &MsgObj{
			Command: "PART",
			Params: []string{
				fmt.Sprintf("#%s", v),
			},
		}
		queue.Add(m)
	}
}

func runQueue() {
	go func() {
		cnt := 0
		for {
			if !queue.IsEmpty() && client != nil {
				m := queue.Get()
				msg := &irc.Message{}
				msg.Command = m.Command
				msg.Params = m.Params

				if m.Command == "JOIN" {
					if indexOf(channels, m.Params[0][1:]) != -1 {
						continue
					}
					channels = append(channels, m.Params[0][1:])
				} else if m.Command == "PART" {
					if indexOf(channels, m.Params[0][1:]) == -1 {
						continue
					}
					idx := indexOf(channels, m.Params[0][1:])
					channels = append(channels[:idx], channels[idx+1:]...)
				}
				fmt.Println("< ", msg.String())
				client.WriteMessage(msg)
			}
			cnt++
			if cnt > 1800 {
				// call rejoin
				ReJoin()
				cnt = 0
			}
			time.Sleep(time.Second * 1)
		}
	}()
}

func ircHandle(c *irc.Client, m *irc.Message) {
	fmt.Println("> ", m.String())
}

func indexOf(c []string, data string) int {
	if len(c) == 0 {
		return -1
	}
	for k, v := range c {
		if v == data {
			return k
		}
	}
	return -1
}
