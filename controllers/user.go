package controllers

import (
	"bluebell/logic"
	"bluebell/models"
	"bluebell/pkg/jwt"
	"errors"

	"bluebell/models/dto"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// SignUp 用户注册
// @Summary 用户注册
// @Description 创建新用户账号
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body dto.SignUpRequest true "注册参数"
// @Success 200 {object} dto.TokenResponse "注册成功"
// @Failure 400 {object} dto.ErrorResponse "请求参数错误"
// @Failure 500 {object} dto.ErrorResponse "服务器内部错误"
// @Router /register [post]
func SignUp(ctx *gin.Context) {
	// 使用 DTO 接受参数
	var req dto.SignUpRequest

	//获取参数验证和参数验证
	if err := ctx.ShouldBind(&req); err != nil {
		zap.L().Error("获取请求参数失败", zap.Error(err))

		//转换为validator.ValidatorErrors类型的error 方便，转为为中文的错误表示
		// ValidationErrors 表示验证错误，而非json格式错误
		errs, ok := err.(validator.ValidationErrors)

		// ok == ture 表示参数错误，而非SignUp
		if ok {
			//ValidatorErros,翻译之后再返回
			if trans != nil {
				Fail(ctx,ErrValidation,errs[0].Translate(trans))
			} else {
				Fail(ctx,ErrValidation,errs.Error())
			}
			return
		} else { // ok == false  表示是JSON格式错误
			FailWithDefault(ctx, ErrInvalidJSON)
			return
		}

	}

	// 转换为内部的 model
	param := &models.ParamSignUp{
		Username:   req.Username,
		Password:   req.Password,
		Repassword: req.RePassword,
	}

	//业务处理, 用户注册
	if err := logic.SignUp(param); err != nil {
		// 判断是否是用户已存在错误
		if errors.Is(err, logic.ErrUserExists) {
			FailWithDefault(ctx, ErrUserExists)
			return
		}
		zap.L().Error("插入用户数据失败", zap.Error(err))
		//返回响应,服务器内部错误
		FailWithDefault(ctx, ErrInternal)
		return
	}

	//返回响应
	Success(ctx, nil)
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录接口，返回 JWT Token
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "登录参数"
// @Success 200 {object} dto.TokenResponse "登录成功"
// @Failure 400 {object} dto.ErrorResponse "请求参数错误"
// @Failure 500 {object} dto.ErrorResponse "服务器内部错误"
// @Router /login [post]
func Login(ctx *gin.Context) {
	// 使用DTO接收请求
	var req dto.LoginRequest

	if err := ctx.ShouldBind(&req); err != nil {
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

	// 转换为内部的 ParamLogin Model
	param := &models.ParamLogin{
		Username: req.Username,
		Password: req.Password,
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
