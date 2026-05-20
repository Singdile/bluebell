package controllers

import (
	"bluebell/logic"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// -----所有和社区相关的-----

// CommunityHandler 查询所有的社区，以(community_id, community_name)的形式返回
// @Summary 获取社区列表
// @Description 获取所有社区的列表信息
// @Tags 社区
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} dto.CommunityListResponse "获取成功"
// @Failure 500 {object} dto.ErrorResponse "服务器内部错误"
// @Router /v1/community [get]
func CommunityHandler(ctx *gin.Context) {
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

// CommunityByID 通过id获取社区详细描述
// @Summary 获取社区详情
// @Description 根据社区ID获取详细信息
// @Tags 社区
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "社区ID"
// @Success 200 {object} dto.CommunityDetailResponse "获取成功"
// @Failure 400 {object} dto.ErrorResponse "请求参数错误"
// @Failure 500 {object} dto.ErrorResponse "服务器内部错误"
// @Router /v1/community/{id} [get]
func CommunityByID(ctx *gin.Context) {
	//获取路径参数id
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)

	if err != nil {
		zap.L().Error("Invalid community_id parameter", zap.String("id", idstr),
			zap.Error(err))
		Fail(ctx, ErrValidation, "无效的社区ID")
		return
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
