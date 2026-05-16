package router

import (
	"bluebell/controllers"
	"bluebell/logger"
	"bluebell/middlewares"
	"bluebell/settings"
	"fmt"

	"time"

	"github.com/swaggo/files"

	_ "bluebell/docs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	v1 := r.Group("/v1", middlewares.JwtAuthMiddleware(),middlewares.RateLimitMiddleware(1,5))

	{
		//获取社区列表
		v1.GET("/community", controllers.CommunityHandler)
		v1.GET("/community/:id", controllers.CommunityByID)

		v1.POST("/post", controllers.PostHandler)
		v1.GET("/post/:id", controllers.GetPostDetailByID)

		// 帖子列表接口(分页)
		//	r.GET("/posts", middlewares.JwtAuthMiddleware(), controllers.GetPostList)
		v1.GET("/posts2", controllers.GetPostListByOrder)
		v1.GET("/post2community", controllers.GetPostListByCommunity)

		// 帖子投票
		v1.POST("/vote", controllers.PostVote)

	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return
}
