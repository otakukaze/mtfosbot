package routes

import (
	"fmt"
	"log"

	"git.trj.tw/golang/mtfosbot/module/config"
	"git.trj.tw/golang/mtfosbot/module/context"
	"git.trj.tw/golang/mtfosbot/router/api"
	"git.trj.tw/golang/mtfosbot/router/google"
	"git.trj.tw/golang/mtfosbot/router/line"
	"git.trj.tw/golang/mtfosbot/router/private"
	"git.trj.tw/golang/mtfosbot/router/rimg"
	"git.trj.tw/golang/mtfosbot/router/twitch"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// NewServ - create new gin server
func NewServ() *gin.Engine {
	r := gin.New()

	conf := config.GetConf()

	redisStr := fmt.Sprintf("%s:%d", conf.Redis.Host, conf.Redis.Port)

	store, err := sessions.NewRedisStore(10, "tcp", redisStr, "", []byte("seckey"))
	if err != nil {
		log.Fatal(err)
	}

	// access log
	r.Use(gin.Logger())
	// error catch
	r.Use(gin.Recovery())
	// enable cors
	r.Use(cors.Default())
	// session
	r.Use(sessions.Sessions("ginsess", store))

	return r
}

// SetRoutes - set routes
func SetRoutes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})

	imageProcGroup := r.Group("/image")
	{
		imageProcGroup.GET("/origin/:imgname", context.PatchCtx(rimg.GetOriginImage))
		imageProcGroup.GET("/thumbnail/:imgname", context.PatchCtx(rimg.GetThumbnailImage))
		imageProcGroup.GET("/line_log_image/:imgname", context.PatchCtx(rimg.GetLineLogImage))
	}

	apiGroup := r.Group("/api")
	{
		apiGroup.POST("/login", context.PatchCtx(api.UserLogin))
		apiGroup.POST("/logout", context.PatchCtx(api.UserLogout))
		apiGroup.GET("/line_msg", context.PatchCtx(api.CheckSession), context.PatchCtx(api.GetLineMessageLog))
		apiGroup.GET("/session", context.PatchCtx(api.CheckSession), context.PatchCtx(api.GetSessionData))
		apiGroup.GET("/twitch/channel/:chid/opay/bar", context.PatchCtx(api.GetDonateBarStatus))
	}

	privateAPIGroup := apiGroup.Group("/private", context.PatchCtx(private.VerifyKey))
	{
		privateAPIGroup.GET("/pages", context.PatchCtx(private.GetFacebookPageIDs))
		privateAPIGroup.POST("/pageposts", context.PatchCtx(private.UpdateFacebookPagePost))
	}

	apiTwitchGroup := apiGroup.Group("/twitch", context.PatchCtx(api.CheckSession))
	{
		apiTwitchGroup.GET("/channels", context.PatchCtx(api.GetChannels), context.PatchCtx(api.GetChannelList))
		twitchChannelGroup := apiTwitchGroup.Group("/channel/:chid", context.PatchCtx(api.GetChannels))
		{
			twitchChannelGroup.GET("/", context.PatchCtx(api.GetChannelData))
			twitchChannelGroup.PUT("/botjoin", context.PatchCtx(api.BotJoinChannel))
			twitchChannelGroup.PUT("/opay", context.PatchCtx(api.OpayIDChange))
			twitchChannelGroup.GET("/opay/setting", context.PatchCtx(api.GetDonateSetting))
			twitchChannelGroup.PUT("/opay/setting", context.PatchCtx(api.UpdateDonateSetting))
		}
	}

	r.POST("/line", context.PatchCtx(line.GetRawBody), context.PatchCtx(line.VerifyLine), context.PatchCtx(line.GetLineMessage))
	lineApis := r.Group("/line")
	{
		lineApis.POST("/", context.PatchCtx(line.GetRawBody), context.PatchCtx(line.VerifyLine), context.PatchCtx(line.GetLineMessage))
	}

	googleApis := r.Group("/google")
	{
		googleApis.GET("/youtube/webhook", context.PatchCtx(google.VerifyWebhook))
		googleApis.POST("/youtube/webhook", context.PatchCtx(google.GetNotifyWebhook))
	}

	twitchApis := r.Group("/twitch")
	{
		twitchApis.GET("/login", context.PatchCtx(twitch.OAuthLogin))
		twitchApis.GET("/oauth", context.PatchCtx(twitch.OAuthProc))
	}

	// set pprof router
	popts := &pprof.Options{
		RoutePrefix: "/dev/pprof",
	}
	pprof.Register(r, popts)
}
