package controllers

import (
	"bluebell/logic"
	"bluebell/models"
	"net/http"
	"strconv"

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
	if err := logic.CreatePost(post); err != nil {
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
	postdetail, err := logic.GetPostDetalByID(post_id)

	if err != nil { //获取失败，返回 错误响应
		zap.L().Error("获取数据库中的postdetail失败", zap.Error(err), zap.Int64("post_id", post_id))
		FailWithDefault(ctx, ErrInternal)
		return

	}

	//成功，返回响应，并携带查询到的数据
	Success(ctx, *postdetail)

}

// GetPostList 获取分页展示的帖子
func GetPostList(ctx *gin.Context) {
	// 获取分页参数
	pagestr := ctx.DefaultQuery("page", "1")
	pagesizestr := ctx.DefaultQuery("pagesize", "20")

	page, err := strconv.ParseInt(pagestr, 10, 64)
	if err != nil || page < 1 {
		page = 1
	}

	pagesize, err := strconv.ParseInt(pagesizestr, 10, 64)
	if err != nil {
		pagesize = 20
	}

	// 业务层
	// 查询post列表
	responsedata, err := logic.GetPostList(page, pagesize)
	if err != nil {
		zap.L().Error("logic.GetPostList() fail", zap.Error(err))
		FailWithDefault(ctx, ErrInternal) //500, 服务器内部错误
		return
	}

	// 返回响应
	Success(ctx, responsedata)

}
