package controllers

import (
	"bluebell/dao/redis"
	"bluebell/logic"
	"bluebell/models"
	"bluebell/models/dto"
	"errors"

	"strconv"

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// PostHandler 创建帖子
// @Summary 创建帖子
// @Description 创建新的帖子
// @Tags 帖子
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.CreatePostRequest true "帖子参数"
// @Success 200 {object} dto.SuccessResponse "创建成功"
// @Failure 400 {object} dto.ErrorResponse "请求参数错误"
// @Failure 401 {object} dto.ErrorResponse "未授权"
// @Failure 500 {object} dto.ErrorResponse "服务器内部错误"
// @Router /v1/post [post]
func PostHandler(ctx *gin.Context) {
	// 通过 DTO 接受请求
	var req dto.CreatePostRequest

	if err := ctx.ShouldBind(&req); err != nil { //gin框架绑定使用了validator库, 参数验证参考了定义结构体的binding tag
		zap.L().Error("获取请求参数失败", zap.Error(err))
		//转换为validator.ValidatorErrors类型的error 方便，转为为中文的错误表示
		// ValidationErrors 表示验证错误，而非json格式错误
		errs, ok := err.(validator.ValidationErrors)

		// ok == ture 表示参数错误，而非PostHandler
		if ok {
			//ValidatorErros,翻译之后再返回
			if trans != nil {
				// validator.ValidationErrors, errs 是一个切片
				Fail(ctx, ErrValidation, errs[0].Translate(trans))
			} else {
				Fail(ctx, ErrValidation, errs[0].Error())
			}
			return
		} else { // ok == false  表示是JSON格式错误
			FailWithDefault(ctx, ErrInvalidJSON)
			return
		}

	}

	// 转换为内部 model
	post := &models.Post{
		Content:     req.Content,
		Title:       req.Title,
		CommunityID: req.CommunityID,
	}

	// 通过jwt认证的用户，会在ctx中保存userid
	user_id, err := getCurrentUser(ctx)
	if err != nil {
		zap.L().Error("用户需要先登录，才能发帖", zap.Error(err))
		FailWithDefault(ctx, ErrLoginFailed)
		return
	}

	post.AuthorID = user_id

	//业务处理
	//创建帖子 Post
	if err := logic.CreatePost(ctx, post); err != nil {
		zap.L().Error("logic.CreatePost() failed", zap.Error(err))
		FailWithDefault(ctx, ErrInternal)
		return
	}

	//返回响应
	Success(ctx, nil)
}

// GetPostDetailByID 根据id,查询对应的post
// @Summary 获取帖子详情
// @Description 根据帖子ID获取详细信息
// @Tags 帖子
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "帖子ID"
// @Success 200 {object} dto.PostDetailResponse "获取成功"
// @Failure 400 {object} dto.ErrorResponse "请求参数错误"
// @Failure 500 {object} dto.ErrorResponse "服务器内部错误"
// @Router /v1/post/{id} [get]
func GetPostDetailByID(ctx *gin.Context) {
	//获取路径参数上的 post_id
	pidstr := ctx.Param("id")

	// 将字符串装换为64位的int
	post_id, err := strconv.ParseInt(pidstr, 10, 64)
	if err != nil {
		zap.L().Error("post_id 路径参数错误", zap.Error(err))
		FailWithDefault(ctx, ErrValidation)
		return
	}

	//业务逻辑
	//查询数据库，获取对应的postdetail
	postdetail, err := logic.GetPostDetailByID(post_id)

	if err != nil { //获取失败，返回 错误响应
		zap.L().Error("获取数据库中的postdetail失败", zap.Error(err), zap.Int64("post_id", post_id))
		FailWithDefault(ctx, ErrInternal)
		return

	}

	//成功，返回响应，并携带查询到的数据
	Success(ctx, *postdetail)

}

// GetPostList 获取分页展示的帖子
// func GetPostList(ctx *gin.Context) {
//	// 获取分页参数
//	pagestr := ctx.DefaultQuery("page", "1")
//	pagesizestr := ctx.DefaultQuery("pagesize", "20")

//	page, err := strconv.ParseInt(pagestr, 10, 64)
//	if err != nil || page < 1 {
//		page = 1
//	}

//	pagesize, err := strconv.ParseInt(pagesizestr, 10, 64)
//	if err != nil {
//		pagesize = 20
//	}

//	// 业务层
//	// 查询post列表
//	responsedata, err := logic.GetPostList(page, pagesize)
//	if err != nil {
//		zap.L().Error("logic.GetPostList() fail", zap.Error(err))
//		FailWithDefault(ctx, ErrInternal) //500, 服务器内部错误
//		return
//	}

//	// 返回响应
//	Success(ctx, responsedata)

// }

// GetPostListByOrder 获取帖子列表（按排序）
// @Summary 获取帖子列表
// @Description 分页获取帖子列表，可按时间或分数排序
// @Tags 帖子
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码" default(1)
// @Param pagesize query int false "每页数量" default(10)
// @Param order query string false "排序方式" Enums(time, score) default(time)
// @Success 200 {object} dto.PostListResponse "获取成功"
// @Failure 400 {object} dto.ErrorResponse "请求参数错误"
// @Failure 500 {object} dto.ErrorResponse "服务器内部错误"
// @Router /v1/posts [get]
func GetPostListByOrder(ctx *gin.Context) {
	// 使用 DTO 接收参数
	var req = dto.PostListRequest{
		Page:        1,
		Pagesize:    20,
		Order:       "time",
		CommunityID: 0,
	}

	//绑定传递的参数,注意部分绑定成功的会改变数据的
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		zap.L().Error(
			"post query 参数绑定失败",
			zap.Error(err),
			zap.Int64("page", req.Page),
			zap.Int64("pagesize", req.Pagesize),
			zap.String("order", req.Order),
		)
	}

	// 转换为内部 model
	postquery := &models.ParamPostQuery{
		Page:        req.Page,
		Pagesize:    req.Pagesize,
		Order:       req.Order,
		CommunityID: req.CommunityID,
	}

	//查询帖子列表
	responsedata, err := logic.GetPostList(ctx, postquery)

	if err != nil {
		zap.L().Error("logic.GetPostListByOrder() fail", zap.Error(err))
		FailWithDefault(ctx, ErrInternal) //500, 服务器内部错误
		return
	}

	//返回响应
	Success(ctx, responsedata)

}

// GetPostListByCommunity 获取指定社区的帖子列表
// @Summary 获取社区帖子列表
// @Description 分页获取指定社区的帖子列表
// @Tags 帖子
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码" default(1)
// @Param pagesize query int false "每页数量" default(10)
// @Param community_id query int true "社区ID"
// @Success 200 {object} dto.PostListResponse "获取成功"
// @Failure 400 {object} dto.ErrorResponse "请求参数错误"
// @Failure 500 {object} dto.ErrorResponse "服务器内部错误"
// @Router /v1/post2community [get]
func GetPostListByCommunity(ctx *gin.Context) {
	// 使用 DTO 接收参数
	var req = dto.PostListRequest{
		Page:        1,
		Pagesize:    20,
		Order:       "time",
		CommunityID: 0,
	}

	// 绑定更新传递的query参数，部分绑定成功的会改变数据
	if err := ctx.ShouldBindQuery(&req); err != nil {
		zap.L().Error(
			"post query 参数绑定失败",
			zap.Error(err),
			zap.Int64("page", req.Page),
			zap.Int64("pagesize", req.Pagesize),
			zap.String("order", req.Order),
			zap.Int64("communityid", req.CommunityID),
		)
	}

	// 使用内部的model
	defaultquery := &models.ParamPostQuery{
		Page:        req.Page,
		Pagesize:    req.Pagesize,
		Order:       req.Order,
		CommunityID: req.CommunityID,
	}

	// 按照参数查询帖子列表
	responsedata, err := logic.GetPostList(ctx, defaultquery)

	if err != nil {
		zap.L().Error("logic.GetPostListByCommunity() failed", zap.Error(err))
		FailWithDefault(ctx, ErrInternal) //500, 服务器内部错误
		return

	}
	// 返回响应
	Success(ctx, responsedata)
}

// 投票
// @Summary 帖子投票
// @Description 对帖子进行赞成、反对或弃票操作
// @Tags 投票
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.VoteRequest true "投票参数"
// @Success 200 {object} dto.SuccessResponse "投票成功"
// @Failure 400 {object} dto.ErrorResponse "请求参数错误"
// @Failure 401 {object} dto.ErrorResponse "未授权"
// @Failure 500 {object} dto.ErrorResponse "服务器内部错误"
// @Router /v1/vote [post]
func PostVote(ctx *gin.Context) {
	// 使用dto来接受请求数据
	var req dto.VoteRequest

	//用户id,帖子id,投票类型
	if err := ctx.ShouldBindJSON(&req); err != nil {
		zap.L().Error("参数绑定错误", zap.Error(err))
		// ValidationErrors 表示参数验证错误，比如投票参数非1,-1,0
		errs, ok := err.(validator.ValidationErrors) //接口类型断言
		if ok {
			Fail(ctx, ErrValidation, fmt.Sprintf("%v", errs.Translate(trans)))
			return
		}

		//josn格式有问题 无法绑定
		FailWithDefault(ctx, ErrInvalidJSON)
		return

	}

	// 转换前端传递的投票情况  up down cancel
	var direction int
	switch req.Action {
	case "up":
		direction = 1
	case "down":
		direction = -1
	case "cancel":
		direction = 0
	}

	// 转换为内部模型
	p := &models.ParamVote {
		PostID: req.PostID,
		Direction:direction,
	}

	user_id, err := getCurrentUser(ctx)
	if err != nil {
		Fail(ctx, ErrLoginFailed, "用户id查询失败,需要登录")
		return
	}

	//业务处理，用户投票
	if err := logic.VotePost(ctx, user_id, p); err != nil {
		zap.L().Error("logic.VotePost() failed", zap.Error(err))

		// 判断是否是投票过期错误
		if errors.Is(err,redis.ErrVoteTimeExpired) {
			Fail(ctx,ErrValidation,"投票时间已过，该帖子已经超过7天了")
			return
		}
		FailWithDefault(ctx, ErrInternal) //500,服务器内部错误
		return
	}

	// 成功，返回响应
	Success(ctx, nil)
}
