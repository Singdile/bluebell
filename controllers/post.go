package controllers

import (
	"bluebell/logic"
	"bluebell/models"
	"net/http"
	"strconv"

	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// PostHandler 创建一个新的post
func PostHandler(ctx *gin.Context) {
	//获取参数
	var post = new(models.Post)
	if err := ctx.ShouldBind(post); err != nil { //gin框架绑定使用了validator库, 参数验证参考了定义结构体的binding tag
		zap.L().Error("获取请求参数失败", zap.Error(err))
		//转换为validator.ValidatorErrors类型的error 方便，转为为中文的错误表示
		// ValidationErrors 表示验证错误，而非json格式错误
		errs, ok := err.(validator.ValidationErrors)

		// ok == ture 表示参数错误，而非PostHandler
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

// GetPostByID 根据id,查询对应的post
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

// GetPostListByOrder 获取postlist,根据指定的order= time/score.
func GetPostListByOrder(ctx *gin.Context) {
	//获取query参数 page pagesize order
	postquery := &models.ParamPostQuery{ //初始化一个默认的参数
		Page:     1,
		Pagesize: 20,
		Order:    "time",
	}

	//绑定传递的参数,注意部分绑定成功的会改变数据的
	err := ctx.ShouldBindQuery(postquery)
	if err != nil {
		zap.L().Error(
			"post query 参数绑定失败",
			zap.Error(err),
			zap.Int64("page", postquery.Page),
			zap.Int64("pagesize", postquery.Pagesize),
			zap.String("order", postquery.Order),
		)
	}

	//查询帖子列表
	responsedata, err := logic.GetPostListByOrder(ctx, postquery)

	if err != nil {
		zap.L().Error("logic.GetPostListByOrder() fail", zap.Error(err))
		FailWithDefault(ctx, ErrInternal) //500, 服务器内部错误
		return
	}

	//返回响应
	Success(ctx, responsedata)

}

// GetPostListByCommunity 获取指定社区的帖子集合，可以指定按时间/热度顺序返回
func GetPostListByCommunity(ctx *gin.Context) {
	// 初始化默认的  page pagesize order communityid
	defaultquery := &models.ParamPostQueryCommunity{
		ParamPostQuery: models.ParamPostQuery{
			Page:     1,
			Pagesize: 20,
			Order:    "time",
		},

		CommunityID: 1,
	}
	// 绑定更新传递的query参数，部分绑定成功的会改变数据
	if err := ctx.ShouldBindQuery(defaultquery); err != nil {
		zap.L().Error(
			"post query 参数绑定失败",
			zap.Error(err),
			zap.Int64("page", defaultquery.Page),
			zap.Int64("pagesize", defaultquery.Pagesize),
			zap.String("order", defaultquery.Order),
			zap.Int64("communityid", defaultquery.CommunityID),
		)
	}

	// 按照参数查询帖子列表
	responsedata, err := logic.GetPostListByCommunity(ctx,defaultquery)

	if err != nil {
		zap.L().Error("logic.GetPostListByCommunity() failed", zap.Error(err))
		FailWithDefault(ctx, ErrInternal) //500, 服务器内部错误
		return

	}
	// 返回响应
	Success(ctx, responsedata)
}

// 投票
func PostVote(ctx *gin.Context) {
	//参数校验
	//用户id,帖子id,投票类型
	p := new(models.ParamVote)
	if err := ctx.ShouldBindJSON(p); err != nil {
		zap.L().Error("参数绑定错误", zap.Error(err))
		// ValidationErrors 表示参数验证错误，比如投票参数非1,-1,0
		errs, ok := err.(validator.ValidationErrors) //接口类型断言
		if !ok {
			Fail(ctx, ErrValidation, fmt.Sprintf("%v", errs.Translate(trans)))
			return
		}

		//josn格式有问题 无法绑定
		return

	}

	user_id, err := getCurrentUser(ctx)
	if err != nil {
		Fail(ctx, ErrLoginFailed, "用户id查询失败,需要登录")
		return
	}

	//业务处理，用户投票
	if err := logic.VotePost(ctx, user_id, p); err != nil {
		zap.L().Error("logic.VotePost() failed", zap.Error(err))
		FailWithDefault(ctx, ErrInternal) //500,服务器内部错误
	}

	// 成功，返回响应
	Success(ctx, nil)
}
