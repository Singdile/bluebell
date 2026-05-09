package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"bluebell/pkg/snowflake"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreatePost 创建post
func CreatePost(ctx *gin.Context,post *models.Post) (err error) {
	// 创建post_id
	post_id := snowflake.GenID()
	post.ID = post_id

	// 初始化CreateTime
	post.CreateTime = time.Now()

	//保存到数据库
	err = mysql.InsertPost(post)
	if err != nil {
		return err
	}
	//保存到reids
	err = redis.CreatePost(ctx, post)
	if err != nil {
		return err
	}
	//返回
	return nil

}

// GetPostByID 根据ID查询post
func GetPostDetalByID(post_id int64) (post *models.PostDetail, err error) {
	// 查询数据库，获取对应的信息
	post, err = mysql.GetPostDetailByID(post_id)
	return
}

// GetPostList 获取用户指定的第几页的帖子数据
func GetPostList(page, pagesize int64) (*models.PostListResponse, error) {
	// 参数校验
	if page < 1 {
		page = 1
	}

	if pagesize < 1 || pagesize > 50 {
		pagesize = 20
	}
	// 调用dao层 查询帖子列表和总数
	list, total, err := mysql.GetPostList(page, pagesize)

	if err != nil {
		zap.L().Error("mysql.GetPostList failed", zap.Error(err))
		return nil, err
	}

	// 计算总页数
	totalpage := int64(total) / pagesize
	if totalpage%pagesize != 0 {
		totalpage++
	}

	//组装响应结构体
	responsedata := &models.PostListResponse{
		Total:      total,
		Page:       page,
		PageSize:   pagesize,
		TotalPages: totalpage,
		List:       list,
	}
	return responsedata, nil
}

// VotePost 用户为帖子投票
func VotePost(ctx *gin.Context,user_id int64,vote * models.ParamVote) error{
	return redis.VotePost(ctx,fmt.Sprintf("%v",user_id),fmt.Sprintf("%v",vote.PostID),float64(vote.Direction))
}
