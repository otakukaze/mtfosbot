package main

import (
	"log"
	"strconv"
	"strings"

	"git.trj.tw/golang/mtfosbot/model"
	"git.trj.tw/golang/mtfosbot/module/background"
	"git.trj.tw/golang/mtfosbot/module/config"
	twitchirc "git.trj.tw/golang/mtfosbot/module/twitch-irc"
	"git.trj.tw/golang/mtfosbot/router/routes"
	"github.com/gin-gonic/gin"
)

var server *gin.Engine

func main() {
	err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	background.SetBackground()

	// create http server
	server = routes.NewServ()
	routes.SetRoutes(server)

	// connect to database
	db, err := model.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = twitchirc.InitIRC()
	if err != nil {
		log.Println(err)
	}

	server.Run(strings.Join([]string{":", strconv.Itoa(config.GetConf().Port)}, ""))
}
