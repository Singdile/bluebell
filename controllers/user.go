package controllers

import (
	"bluebell/logic"
	"bluebell/models"
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
		if trans != nil {
			ctx.JSON(http.StatusOK, gin.H{"msg": errs.Translate(trans)})
		} else {
			ctx.JSON(http.StatusOK, gin.H{"msg": errs.Error()})
		}
		return
	}

	//业务处理, 用户注册
	if err := logic.SignUp(param); err != nil {
		zap.L().Error("插入用户数据失败", zap.Error(err))
		//返回响应,服务器内部错误
		FailWithDefault(ctx, ErrInternal)
		return
	}

	//返回响应
	Success(ctx, nil)
}

// Login 登陆
func Login(ctx *gin.Context) {
	//获取登录用户信息
	param := new(models.ParamLogin)

	if err := ctx.ShouldBind(param); err != nil {
		zap.L().Error("参数绑定失败", zap.Error(err))

		// 断言是否是validator的验证错误
		_, ok := err.(*validator.InvalidValidationError)

		if !ok {
			// validator解析错误
			FailWithDefault(ctx, ErrInvalidJSON)
			return
		}

		// validator 的验证错误， 比如某个参数为空了/ 或者两个参数不等
		FailWithDefault(ctx, ErrValidation)
		return
	}

	// 业务逻辑处理
	// 验证用户信息是否合法
	if err := logic.Login(param); err != nil {
		zap.L().Error("验证失败", zap.Error(err))
		FailWithDefault(ctx, ErrValidation)
		return
	}

	// 登录成功
	Success(ctx, nil)
}
