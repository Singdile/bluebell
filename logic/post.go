package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"bluebell/pkg/snowflake"
	"fmt"

	"time"

	"github.com/gin-gonic/gin"
)

// CreatePost 创建post
func CreatePost(ctx *gin.Context, post *models.Post) (err error) {
	// 创建post_id
	post_id := snowflake.GenID()
	post.ID = post_id

	// 初始化CreateTime
	post.CreateTime = time.Now()

	// 初始化 Score（等于创建时间戳）
	post.Score = post.CreateTime.Unix()

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

// GetPostListByOrder
// 先查reids
// reids为空，则使用mysql
func getPostListByOrder(ctx *gin.Context, postquery *models.ParamPostQuery) (*models.PostListResponse, error) {
	// 先获取帖子的总数
	total, err := getTotalPostCount(ctx)
	if err != nil {
		return nil, err
	}

	// 计算总页数
	totalpage := total / postquery.Pagesize
	if total%postquery.Pagesize != 0 {
		totalpage++
	}

	// 分页查询数据
	//1.尝试从redis按照score/time 的降序取出 post_id
	postids, err := redis.GetPostIDs(ctx, postquery)

	// redis 正常有数据 ——> 返回redis结果
	if err == nil && len(postids) > 0 {
		responsedata, err := buildResponseFromRedis(ctx, postids, postquery)
		if err != nil {
			return nil, err
		}
		responsedata.Total = total
		responsedata.TotalPages = totalpage
		return responsedata, nil
	}

	//2. redis 查询失败 ——> 查询mysql
	responsedata, err := getPostListFromMySQL(ctx, postquery)
	if err != nil {
		return nil, err
	} else {
		responsedata.Total = total
		responsedata.TotalPages = totalpage
		return responsedata, nil
	}

}

// buildResponseFromRedis redis查询到postids,再到mysql中查找对应的
func buildResponseFromRedis(ctx *gin.Context, postids []string, postquery *models.ParamPostQuery) (*models.PostListResponse, error) {
	//按照post_id到mysql数据库中查询post, 返回的数据的顺序要是postids中的顺序
	list, err := mysql.GetPostListByIDs(postids, 0)
	if err != nil {
		return nil, err
	}

	//构造成响应结构体返回
	responsedata := &models.PostListResponse{
		Page:     postquery.Page,
		PageSize: postquery.Pagesize,
		List:     list,
	}

	return responsedata, nil

}

// getPostListFromMySQL 从mysql中分页查询postlist
func getPostListFromMySQL(ctx *gin.Context, postquery *models.ParamPostQuery) (*models.PostListResponse, error) {
	// MySQL 按序分页查询获取数据
	list, err := mysql.GetPostListByOrder(postquery.Page, postquery.Pagesize, postquery.Order)
	if err != nil {
		return nil, err
	}

	// 构造返回
	responsedata := &models.PostListResponse{
		Page:     postquery.Page,
		PageSize: postquery.Pagesize,
		List:     list,
	}

	return responsedata, nil

}

// getPostListByCommunity 获取社区帖子列表
func getPostListByCommunity(ctx *gin.Context, query *models.ParamPostQuery) (*models.PostListResponse, error) {
	// 获取社区帖子总数
	total, err := getCommunityPostCount(ctx, query.CommunityID)
	if err != nil {
		return nil, err
	}

	// 计算总的记录数total 和 总页数 totalpage
	totalpage := total / query.Pagesize
	if total%query.Pagesize != 0 {
		totalpage++
	}

	// 查询redis
	if total == 0 { // 该社区还没有帖子
		return &models.PostListResponse{
			Total: 0,
			Page: query.Page,
			PageSize: query.Pagesize,
			TotalPages: 0,
			List: []*models.PostListItem{},
		},nil
	}


	// 查询社区按照对应的order的post_id 列表
	postids, err := redis.GetPostIDs(ctx, query)
	if err != nil {
		return nil, err
	}

	// 按照post_id到mysql数据库中查询数据，返回的数据顺序是post_id 列表的顺序
	list, err := mysql.GetPostListByIDs(postids, query.CommunityID)
	if err != nil {
		return nil, err
	}

	// 构造响应结构体返回
	responsedata := &models.PostListResponse{
		Total:      total,
		Page:       query.Page,
		PageSize:   query.Pagesize,
		TotalPages: totalpage,
		List:       list,
	}
	return responsedata, nil
}

// GetPostList 统一查询帖子列表的接口(全局，社区查询)
func GetPostList(ctx *gin.Context, query *models.ParamPostQuery) (*models.PostListResponse, error) {
	if query.CommunityID == 0 {
		return getPostListByOrder(ctx, query)
	} else {
		return getPostListByCommunity(ctx, query)
	}
}

// VotePost 用户为帖子投票
func VotePost(ctx *gin.Context, user_id int64, vote *models.ParamVote) error {
	return redis.VotePost(ctx, fmt.Sprintf("%v", user_id), fmt.Sprintf("%v", vote.PostID), float64(vote.Direction))
}

// getTotalPostCount 获取帖子总数(先redis， fallback mysql)
func getTotalPostCount(ctx *gin.Context) (int64, error) {
	//优先从redis中获取
	count, err := redis.GetTotalPostCount(ctx)
	if err == nil && count > 0 {
		return count, nil
	}

	// reids 获取失败，fallback 到 mysql
	return mysql.GetTotalPostCount()
}

// getCommunityPostCount 获取社区帖子总数 (先redis， fallback mysql)
func getCommunityPostCount(ctx *gin.Context, communityid int64) (int64, error) {
	// 优先从redis中获取
	count, err := redis.GetCommunityPostCount(ctx, communityid)
	if err == nil && count > 0 {
		return count, nil
	}

	// redis 获取失败， fallbcak 到 mysql
	return mysql.GetCommunityPostCount(communityid)
}
