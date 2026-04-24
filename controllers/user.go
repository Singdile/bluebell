package controllers

import (
	"bluebell/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// 用户注册
func SignUp(ctx *gin.Context) {
	param := new(models.ParamSignUp)
	//获取参数验证和参数验证
	if err := ctx.ShouldBind(param); err != nil {
		zap.L().Error("获取请求参数失败", zap.Error(err))

		//转换为validator.ValidatorErrors类型的error 方便，转为为中文的错误表示
		errs, ok := err.(validator.ValidationErrors)

		//非ValidatorErrors,直接返回
		if !ok {
			ctx.JSON(http.StatusOK, gin.H{"msg": "请求参数有误", "err": err.Error()})
			return
		}

		//ValidatorErros,翻译之后再返回
		ctx.JSON(http.StatusOK, gin.H{"msg": errs.Translate(trans)})
		return
	}

	fmt.Printf("user: %v\n", param)
	//用户注册
	// logic.SignUp()
	//返回响应
	ctx.JSON(http.StatusOK, gin.H{"msg": "success"})

}
