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
func CreatePost(ctx *gin.Context, post *models.Post) (err error) {
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
	return err

}

// GetPostByID 根据ID查询post
func GetPostDetailByID(post_id int64) (*models.PostDetail, error) {
	// 查询数据库，获取对应的信息
	return mysql.GetPostDetailByID(post_id)
}

// GetPostList 获取用户指定的第几页的帖子数据
// func GetPostList(page, pagesize int64) (*models.PostListResponse, error) {
//	// 参数校验
//	if page < 1 {
//		page = 1
//	}

//	if pagesize < 1 || pagesize > 50 {
//		pagesize = 20
//	}
//	// 调用dao层 查询帖子列表和总数
//	list, total, err := mysql.GetPostList(page, pagesize)

//	if err != nil {
//		zap.L().Error("mysql.GetPostList failed", zap.Error(err))
//		return nil, err
//	}

//	// 计算总页数
//	totalpage := int64(total) / pagesize
//	if total%pagesize != 0 {
//		totalpage++
//	}

//	//组装响应结构体
//	responsedata := &models.PostListResponse{
//		Total:      total,
//		Page:       page,
//		PageSize:   pagesize,
//		TotalPages: totalpage,
//		List:       list,
//	}
//	return responsedata, nil
// }

// GetPostListByOrder
func GetPostListByOrder(ctx *gin.Context, postquery *models.ParamPostQuery) (*models.PostListResponse, error) {
	//1.从redis按照score/time 的降序取出 post_id
	postids, err := redis.GetPostIDsInOrder(ctx, postquery)
	if err != nil {
		return nil, err
	}

	if len(postids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder() success,bu get 0 record")
		return nil, nil
	}
	//2.按照post_id到mysql数据库中查询post, 返回的数据的顺序要是postids中的顺序
	list, total, err := mysql.GetPostListByIDs(postids,0)
	if err != nil {
		return nil, err
	}

	//查询每篇帖子的投票数
	votePList, voteNList, err := redis.GetPostVoteData(ctx, postids)
	if err != nil {
		return nil, err
	}

	// 计算总页数
	totalpage := int64(total) / postquery.Pagesize
	if total%postquery.Pagesize != 0 {
		totalpage++
	}

	//3.构造成响应结构体返回
	responsedata := &models.PostListResponse{
		Total:      total,
		Page:       postquery.Page,
		PageSize:   postquery.Pagesize,
		TotalPages: totalpage,
		List:       list,
		VoteP:      votePList,
		VoteN:      voteNList,
	}
	return responsedata, nil
}

func GetPostListByCommunity(ctx *gin.Context, query *models.ParamPostQueryCommunity) (*models.PostListResponse, error) {
	// 查询社区按照对应的order的post_id 列表
	postids, err := redis.GetPostIDsInOrderByCom(ctx, query)
	if err != nil {
		return nil, err
	}

	if len(postids) == 0 {
		//该社区没有帖子,返回空列表
		return &models.PostListResponse {
			Total: 0,
			Page: query.Page,
			PageSize: query.Pagesize,
			TotalPages: 0,
			List: []*models.PostListItem{},
			VoteP: []string{},
			VoteN: []string{},
		}, nil
	}
	// 按照post_id到mysql数据库中查询数据，返回的数据顺序是post_id 列表的顺序
	list, total, err := mysql.GetPostListByIDs(postids,query.CommunityID)
	if err != nil {
		return nil, err
	}

	// 根据post_id 列表 查询每篇帖子的投票数
	votePList, voteNList, err := redis.GetPostVoteData(ctx, postids)
	if err != nil {
		return nil, err
	}

	// 计算总页数
	totalpage := int64(len(postids)) / query.Pagesize
	if total%query.Pagesize != 0 {
		totalpage++
	}

	// 构造响应结构体返回
	//3.构造成响应结构体返回
	responsedata := &models.PostListResponse{
		Total:      int64(len(postids)),
		Page:       query.Page,
		PageSize:   query.Pagesize,
		TotalPages: totalpage,
		List:       list,
		VoteP:      votePList,
		VoteN:      voteNList,
	}
	return responsedata, nil
}

// VotePost 用户为帖子投票
func VotePost(ctx *gin.Context, user_id int64, vote *models.ParamVote) error {
	return redis.VotePost(ctx, fmt.Sprintf("%v", user_id), fmt.Sprintf("%v", vote.PostID), float64(vote.Direction))
}
