package controllers

import (
	"bluebell/logic"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// -----所有和社区相关的-----

// CommunityHandler 查询所有的社区，以(community_id, community_name)的形式返回
func CommunityHandler(ctx *gin.Context) {
	// 参数处理

	// 业务处理
	communitylist, err := logic.GetCommunityList() //获取社区信息列表

	if err != nil {
		zap.L().Error("logic.Getcommunitylist() failed", zap.Error(err))
		FailWithDefault(ctx, ErrInternal) //500,服务器内部错误
		return
	}

	// 返回响应
	Success(ctx, communitylist)

}

// CommunityByID 通过id获取社区详细描述，
func CommunityByID(ctx *gin.Context) {
	//获取路径参数id
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)

	if err != nil {
		zap.L().Error("Invalid community_id parameter", zap.String("id", idstr),
			zap.Error(err))
		Fail(ctx, ErrValidation, "无效的社区ID")
	}

	//业务处理
	communitydetail, err := logic.GetCommunityByID(id)

	//检查数据库错误
	if err != nil {
		zap.L().Error("logic.GetCommunityByID() failed", zap.Int64("community_id", id), zap.Error(err))
		FailWithDefault(ctx, ErrIdNotExist) // 500 ,internal error
		return
	}
	//返回响应
	Success(ctx, communitydetail)
}
