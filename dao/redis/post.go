package redis

import (
	"bluebell/models"

	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)



// GetPostIDs 获取指定范围，排序后的帖子 post_id 列表
func GetPostIDs(ctx *gin.Context, postquery *models.ParamPostQuery) ([]string, error){
	if postquery.CommunityID == 0 {
		return getPostIDsInOrder(ctx, postquery)
	} else {
		return getPostIDsInOrderByCom(ctx,postquery)
	}
}


// getPostIDsInOrder 获取全局指定范围内的帖子
func getPostIDsInOrder(ctx *gin.Context, postquery *models.ParamPostQuery) ([]string, error) {
	//判断order
	key := getOrderKey(postquery.Order)
	//查询redis  时间/分数 降序查询
	return getIDsFromKey(ctx, rdb, key, postquery.Page, postquery.Pagesize).Result()

}

// getPostIDsInOrderByCom 取出community下对应order的post_id 列表
func getPostIDsInOrderByCom(ctx *gin.Context, query *models.ParamPostQuery) ([]string, error) {
	//判断排序order,选取对应的key
	orderkey := getOrderKey(query.Order)

	// redis查询key
	// 社区key  bluebell:post:coummunity:community_id
	ckey := KeyPostCommunitySetPrefix + strconv.FormatInt(query.CommunityID, 10)

	// 临时结果key  bluebell:post:community:score/time:community_id
	deskey := KeyCommunityZsetPF + query.Order + strconv.FormatInt(query.CommunityID, 10)

	// 检查 key 对应的 zset是否存在，存在直接查询;不存在则构造
	if  rdb.Exists(ctx, deskey).Val() < 1 {//无缓存
		pipeline := rdb.TxPipeline()
		// communityset 和 KeyPostScoreZset 进行交集操作
		pipeline.ZInterStore(ctx, deskey, &redis.ZStore{
			Keys:      []string{ckey, orderkey},
			Aggregate: "MAX",
		})

		pipeline.Expire(ctx, deskey, 60*time.Second)

		//查询 zset 获取按照规则降序排序的 post_id
		rangeCmd := getIDsFromKey(ctx, pipeline, deskey, query.Page, query.Pagesize)

		// pipeline执行查询
		_, err := pipeline.Exec(ctx)
		if err != nil {
			zap.L().Error("redis pipeline 执行失败", zap.Error(err))
			return nil, err
		}

		// 获取post_id 结果
		postIDs, err := rangeCmd.Result()
		if err != nil {
			zap.L().Error("ZRangeArgs 查询失败", zap.Error(err))
			return nil, err
		}

		return postIDs, nil
	}

	// 并集key存在，直接查询
	rangeCmd := getIDsFromKey(ctx, rdb, deskey, query.Page, query.Pagesize)
	postIDs, err := rangeCmd.Result()
	if err != nil {
		zap.L().Error("ZRangeArgs 查询失败", zap.Error(err))
		return nil, err

	}
	return postIDs, nil

}

// getOrderKey 管理order 到 key 的映射
func getOrderKey(order string) string {
	if order == models.Scoreorder {
		return KeyPostScoreZset
	}

	return KeyPostTimeZset
}

// getPageScope 获取指定范围内的ZSet 数据
func getPageScope(page, pagesize int64) (start, end int64) {
	start = (page - 1) * pagesize
	end = start + pagesize - 1
	return
}

// getIDsFromKey 通用的ZRange查询
func getIDsFromKey(ctx *gin.Context, rdb redis.Cmdable, key string, page, pagesize int64) *redis.StringSliceCmd {
	start, end := getPageScope(page, pagesize)
	return rdb.ZRangeArgs(ctx, redis.ZRangeArgs{
		Key:   key,
		Start: start,
		Stop:  end,
		Rev:   true,
	})
}

// GetPostVoteData 根据ids，查询帖子的投票情况
func GetPostVoteData(ctx *gin.Context, ids []string) ([]string, []string, error) {
	// 帖子同意的票数 和 不同意的票数
	votePList := make([]string, 0, len(ids))
	voteNList := make([]string, 0, len(ids))

	// 一次rtt查询
	pipeline := rdb.TxPipeline()

	for _, id := range ids {
		pipeline.ZCount(ctx, KeyPostVotedZsetPrefix+id, "1", "1")
		pipeline.ZCount(ctx, KeyPostVotedZsetPrefix+id, "-1", "-1")
	}

	//执行查询
	cmders, err := pipeline.Exec(ctx)

	if err != nil {
		zap.L().Error("redis pipeline,Exec failed", zap.Error(err))
		return nil, nil, err

	}

	// 解析结果
	for i := 0; i < len(cmders); i += 2 {
		votePCount := cmders[i].(*redis.IntCmd).Val()
		voteNCount := cmders[i+1].(*redis.IntCmd).Val()

		//转为 string存放返回
		votePList = append(votePList, strconv.FormatInt(votePCount, 10))
		voteNList = append(voteNList, strconv.FormatInt(voteNCount, 10))

	}

	return votePList, voteNList, nil
}
