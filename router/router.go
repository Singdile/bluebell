package router

import (
	"bluebell/controllers"
	"bluebell/logger"
	"bluebell/middlewares"
	"bluebell/settings"
	"fmt"

	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() (r *gin.Engine) {
	r = gin.New()

	//定义validator的翻译选项,当接受到zh 时能够将返回的validator的错误转换为中文
	if err := controllers.InitTrans("zh"); err != nil {
		panic(fmt.Sprintf("init trans failed, err: %v \n", err))
	}

	//  CORS跨域配置
	corsconfig := settings.Conf.Crosconfig

	r.Use(cors.New(cors.Config{
		AllowOrigins: corsconfig.Allow_origins,
		AllowMethods: corsconfig.Allow_methods,
		AllowHeaders: corsconfig.Allow_headers,
		MaxAge:       time.Duration(corsconfig.Max_age) * time.Second,
	}))
	//使用zap接管gin的日志记录
	r.Use(logger.ZapLogger(), logger.ZapRecovery(true))

	//注册路由
	//用户注册
	r.POST("/register", controllers.SignUp)

	r.POST("/login", controllers.Login)

	r.POST("/ping", middlewares.JwtAuthMiddleware(), func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "success auth")
	})

	//获取社区列表
	r.GET("/community", controllers.CommunityHandler)
	r.GET("/community/:id", controllers.CommunityByID)

	r.POST("/post", middlewares.JwtAuthMiddleware(), controllers.PostHandler)
	r.GET("/post/:id", middlewares.JwtAuthMiddleware(), controllers.GetPostDetailByID)

	// 帖子列表接口(分页)
	r.GET("/posts", middlewares.JwtAuthMiddleware(), controllers.GetPostList)
	//r.GET("/community/:id/posts", controllers.GetPostListByCommunity)


	r.POST("/vote",middlewares.JwtAuthMiddleware() ,controllers.PostVote)
	return
}
