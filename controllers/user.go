package controllers

import (
	"bluebell/logic"
	"bluebell/models"
	"bluebell/pkg/jwt"
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
		// ValidationErrors 表示验证错误，而非json格式错误
		errs, ok := err.(validator.ValidationErrors)

		// ok == ture 表示参数错误，而非SignUp
		if ok {
			//ValidatorErros,翻译之后再返回
			if trans != nil {
				ctx.JSON(http.StatusOK, gin.H{"msg": errs.Translate(trans)})
			} else {
				ctx.JSON(http.StatusOK, gin.H{"msg": errs.Error()})
			}
			return
		} else { // ok == false  表示是JSON格式错误
			FailWithDefault(ctx, ErrInvalidJSON)
			return
		}

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
		//  InvalidValidationError - 表示传给validator的参数本身有问题
		//  ValidationErrors - 字段验证失败的错误集合
		_, ok := err.(validator.ValidationErrors)

		if ok { //ok == true , 表示是 ValidatonErros错误
			FailWithDefault(ctx, ErrValidation)
		} else { //表示不合法的JSON格式
			FailWithDefault(ctx, ErrInvalidJSON)
		}

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
	// 生成 Token
	tokenstring, err := jwt.GenToken(param.UserID, param.Username)
	if err != nil {
		zap.L().Error("生成Token失败", zap.Error(err))
		FailWithDefault(ctx, ErrInvalidToken)
		return
	}

	//返回响应，携带tokenstring
	Success(ctx, tokenstring)
}

