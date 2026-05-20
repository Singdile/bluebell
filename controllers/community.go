package controllers

import (
	"bluebell/logic"
	"bluebell/models"
	"bluebell/models/dto"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// -----所有和社区相关的-----

// CommunityHandler 查询所有的社区，以(id, community_name)的形式返回
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
		zap.L().Error("Invalid communityid parameter", zap.String("id", idstr),
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


// CreateCommunity 创建社区
func CreateCommunity(ctx *gin.Context) {
	// 利用dto 获取参数
	var req dto.CreateCommunityRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		Fail(ctx,ErrValidation,"参数验证失败")
		return
	}

	// 构造内部所需要的模型
	userid,err := getCurrentUser(ctx)
	if err != nil {
		Fail(ctx, ErrLoginFailed, "用户未登录")
		return
	}

	community := &models.Community {
		Name: req.CommunityName,
		Introduction: req.Introduction,
		CreatorID: userid,
		Status: 1,
	}
	// 业务处理
	var id int64
	if id,err = logic.CreateCommunity(community); err != nil {
		zap.L().Error("Failed to create community in mysql",zap.Error(err),zap.Any("community info",community))
		Fail(ctx,ErrInternal,err.Error())
		return
	}


	// 返回响应
	Success(ctx, gin.H {
		"community_id": id,
	})
	return
}

// UpdateCommunity 更新社区数据，这里主要是社区的名称和描述信息
func UpdateCommunity(ctx *gin.Context) {
	// 获取请求参数
	var req dto.UpdateCommunityRequest

	if err:= ctx.ShouldBindJSON(&req); err != nil {
		Fail(ctx,ErrValidation,"参数验证错误")
		return
	}

	// 转换为内部的模型
	community := &models.Community {
		ID: req.ID,
		Name: req.CommunityName,
		Introduction: req.Introduction,
	}

	// 业务处理
	if err := logic.UpdateCommunity(community); err != nil {
		Fail(ctx,ErrInternal,err.Error())
		return
	}

	// 返回响应
	Success(ctx,nil)
	return
}
