package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"git.trj.tw/golang/mtfosbot/model"
	"git.trj.tw/golang/mtfosbot/router/routes"
	"github.com/gin-gonic/gin"
)

var server *gin.Engine

func main() {
	portNum, err := strconv.ParseUint(os.Getenv("PORT"), 10, 32)
	if err != nil || portNum < 1024 || portNum > 65535 {
		portNum = 10230
	}

	// create http server
	server = routes.NewServ()
	routes.SetRoutes(server)

	// connect to database
	_, err = model.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	server.Run(strings.Join([]string{":", strconv.Itoa(int(portNum))}, ""))
}
