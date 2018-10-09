package main

import (
	"encoding/gob"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"

	"git.trj.tw/golang/mtfosbot/module/apis/twitch"
	"git.trj.tw/golang/mtfosbot/module/cmd"
	"git.trj.tw/golang/mtfosbot/module/es"
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

	// connect to database
	db, err := model.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if runOptions.DBTool {
		cmd.DBTool()
	}

	registerTypes()
	background.SetBackground()

	// create http server
	server = routes.NewServ()
	routes.SetRoutes(server)

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
	if !utils.CheckExists(conf.LogImageRoot, true) {
		log.Fatal(errors.New("log image root not exists"))
	}

	go func() {
		for {
			PrintMemUsage()
			time.Sleep(time.Second * 20)
		}
	}()

	server.Run(strings.Join([]string{":", strconv.Itoa(config.GetConf().Port)}, ""))
}

// PrintMemUsage -
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
	fmt.Printf("HeapAlloc = %v MiB", bToMb(m.HeapAlloc))
	fmt.Printf("\t HeapSys = %v MiB", bToMb(m.HeapSys))
	fmt.Printf("\t NextGC = %v MiB\n", bToMb(m.NextGC))

	obj := map[string]interface{}{
		"Alloc":        fmt.Sprintf("%v MiB", bToMb(m.Alloc)),
		"Sys":          fmt.Sprintf("%v MiB", bToMb(m.Sys)),
		"HeapAlloc":    fmt.Sprintf("%v MiB", bToMb(m.HeapAlloc)),
		"HeapSys":      fmt.Sprintf("%v MiB", bToMb(m.HeapSys)),
		"HeapIdle":     fmt.Sprintf("%v MiB", bToMb(m.HeapIdle)),
		"HeapInuse":    fmt.Sprintf("%v MiB", bToMb(m.HeapInuse)),
		"HeapReleased": fmt.Sprintf("%v MiB", bToMb(m.HeapReleased)),
		"StackInuse":   fmt.Sprintf("%v MiB", bToMb(m.StackInuse)),
		"StackSys":     fmt.Sprintf("%v MiB", bToMb(m.StackSys)),
		"GCSys":        fmt.Sprintf("%v MiB", bToMb(m.GCSys)),
		"NextGC":       fmt.Sprintf("%v MiB", bToMb(m.NextGC)),
		"NumGC":        fmt.Sprintf("%v", m.NumGC),
	}

	es.PutLog("memory", obj)
}
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
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
	gob.Register(twitch.TwitchTokenData{})
	gob.Register(twitch.UserInfo{})
	gob.Register(map[string]interface{}{})
}
