package routes

import (
	"log"

	"git.trj.tw/golang/mtfosbot/module/context"
	"git.trj.tw/golang/mtfosbot/router/api"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// NewServ - create new gin server
func NewServ() *gin.Engine {
	r := gin.New()

	store, err := sessions.NewRedisStore(10, "tcp", "localhost:6379", "")
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

	apiGroup := r.Group("/api")
	{
		apiGroup.POST("/login", context.PatchCtx(api.UserLogin))
	}
}
