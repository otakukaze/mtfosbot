package twitchirc

import (
	"fmt"
	"net"
	"time"

	"gopkg.in/irc.v2"

	"git.trj.tw/golang/mtfosbot/module/config"
)

var client *irc.Client
var queue *QueueList
var channels []string

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

	queue = NewQueue()
	go runQueue()

	channels = make([]string, 0)
	return
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
			fmt.Sprintf(":%s", msg),
		},
	}
	queue.Add(m)
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

func runQueue() {
	for {
		if !queue.IsEmpty() {
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

		time.Sleep(time.Microsecond * 1500)
	}
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
