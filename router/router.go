package router

import (
	"net/http"
	"web_app/logger"

	"github.com/gin-gonic/gin"
)

func SetupRouter() (r *gin.Engine) {
	r = gin.New()

	//使用zap接管gin的日志记录
	r.Use(logger.ZapLogger(), logger.ZapRecovery(false))

	//注册路由
	r.GET("/", func(ctx *gin.Context) { ctx.String(http.StatusOK, "gin welcome") })

	return
}
