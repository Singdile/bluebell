package controllers

import (
	"errors"

	"github.com/gin-gonic/gin"
)

const Ctxuserid = "userID"

// getCurrentUser 获取当前登录用户的id
func getCurrentUser(ctx *gin.Context) (userid int64, err error) {
	uid, ok := ctx.Get(Ctxuserid)
	if !ok { //获取失败
		//返回响应，用户未登录
		err := errors.New(getErrorMsg(ErrUnauthorized))
		FailWithDefault(ctx, ErrUnauthorized)
		return userid, err
	}

	userid, ok = uid.(int64)

	if !ok {
		//无法转换成功，返回响应
		err := errors.New(getErrorMsg(ErrUnauthorized))
		FailWithDefault(ctx, ErrUnauthorized)
		return userid, err
	}

	return userid, nil
}
