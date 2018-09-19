package main

import (
	"encoding/gob"
	"errors"
	"flag"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"git.trj.tw/golang/mtfosbot/module/options"
	"git.trj.tw/golang/mtfosbot/module/utils"

	"git.trj.tw/golang/mtfosbot/model"
	"git.trj.tw/golang/mtfosbot/module/background"
	"git.trj.tw/golang/mtfosbot/module/config"
	twitchirc "git.trj.tw/golang/mtfosbot/module/twitch-irc"
	"git.trj.tw/golang/mtfosbot/router/routes"
	"github.com/gin-gonic/gin"
)

var server *gin.Engine

func init() {
	options.RegFlag()
	flag.Parse()
}

func main() {
	runOptions := options.GetFlag()

	if runOptions.Help {
		flag.Usage()
		return
	}

	err := config.LoadConfig(runOptions.Config)
	if err != nil {
		log.Fatal(err)
	}

	registerTypes()
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

	go twitchirc.InitIRC()

	// create thumbnail directory
	conf := config.GetConf()
	if !utils.CheckExists(conf.ImageRoot, true) {
		log.Fatal(errors.New("image root not exists"))
	}
	if !utils.CheckExists(path.Join(conf.ImageRoot, "thumbnail"), true) {
		err = os.MkdirAll(path.Join(conf.ImageRoot, "thumbnail"), 0775)
		if err != nil {
			log.Fatal(err)
		}
	}

	server.Run(strings.Join([]string{":", strconv.Itoa(config.GetConf().Port)}, ""))
}

func registerTypes() {
	gob.Register(model.Account{})
	gob.Register(model.Commands{})
	gob.Register(model.DonateSetting{})
	gob.Register(model.FacebookPage{})
	gob.Register(model.KeyCommands{})
	gob.Register(model.LineGroup{})
	gob.Register(model.LineMessageLog{})
	gob.Register(model.LineUser{})
	gob.Register(model.OpayDonateList{})
	gob.Register(model.TwitchChannel{})
	gob.Register(model.YoutubeChannel{})
	gob.Register(map[string]interface{}{})
}
