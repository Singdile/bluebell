package router

import (
	"bluebell/controllers"
	"bluebell/logger"
	"fmt"

	"github.com/gin-gonic/gin"
)

func SetupRouter() (r *gin.Engine) {
	r = gin.New()

	//定义validator的翻译选项,当接受到zh 时能够将返回的validator的错误转换为中文
	if err := controllers.InitTrans("zh"); err != nil {
		panic(fmt.Sprintf("init trans failed, err: %v \n", err))
	}

	//使用zap接管gin的日志记录
	r.Use(logger.ZapLogger(), logger.ZapRecovery(true))

	//注册路由
	//用户注册
	r.POST("/register", controllers.SignUp)

	return
}
