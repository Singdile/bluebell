package redis

import (
	"bluebell/models"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const oneWeekInSeconds = float64(24 * 3600 * 7)
const ErrVoteTimeExpired = "投票时间已经过了"

// 投票功能实现：
// 1.用户投票 (用户向帖子投票)
// 2.

// 投票算法
// 分数 = （赞成 - 不赞成)* 432 + timestamp
// 432的由来， 当我们希望200个赞能让该帖子多停留一天
// 24*60*60 / 200 = 432

/*
投票情况:
direction = 1
-之前没投过，现在投赞成 +432
-之前投不赞成，现在改投赞成 +432*2

direction = 0
-之前投赞成，现在取消 -432
-之前投不赞成，现在取消 +432

direction = -1
-之前没投过，现在投不赞成 -432
-之前投赞成，现在改投不赞成 -432*2

投票的限制：
-给投票的时间设置为一个时间段(7天)，记录该时间段的帖子的投票情况，包括赞成，返回，用户
-过了该时间段之后，不能继续投票。记录下，帖子的总的赞成，不赞成票数
-删除与帖子关联的用户信息
*/

// TODO 完成投票
func VotePost(ctx *gin.Context, user_id, post_id string, direction float64) error {
	//判断投票的限制  超时n不能投票
	posttime := rdb.ZScore(ctx, KeyPostTimeZset, post_id).Val()

	if float64(time.Now().Unix()) > oneWeekInSeconds+posttime {
		return errors.New(ErrVoteTimeExpired)
	}

	//获取用户之前的投票情况,如果没有投过，返回的是零值
	origin_vote := rdb.ZScore(ctx, KeyPostVotedZsetPrefix+post_id, user_id).Val()
	if math.Abs(direction-origin_vote) < 0.0001 { //投票动作没变
		return nil
	}

	//计算需要更新的分数
	var changescore = (direction - origin_vote) * 432

	//更新redis数据
	// 1.开启pipeline事务
	pipe := rdb.TxPipeline()

	// 2.更新post的score
	pipe.ZIncrBy(ctx, KeyPostScoreZset, changescore, post_id)

	// 3.更新用户投票状态
	if direction == 0 { //不投票，移除数据
		pipe.ZRem(ctx, KeyPostVotedZsetPrefix+post_id, user_id)
	} else {
		//记录新状态
		pipe.ZAdd(ctx, KeyPostVotedZsetPrefix+post_id, redis.Z{
			Score:  direction,
			Member: user_id,
		})
	}

	// 4.提交执行
	_, err := pipe.Exec(ctx)
	return err
}

// CreatePost 同步记录创建的post
func CreatePost(ctx *gin.Context, post *models.Post) error {
	//开启pipeline事务操作，pipe负责一次发送，事务负责原子化操作
	pipe := rdb.TxPipeline()

	//记录创建的post_id, create_time
	pipe.ZAdd(ctx, KeyPostTimeZset, redis.Z{
		Score:  float64(post.CreateTime.Unix()),
		Member: fmt.Sprintf("%v", post.ID),
	})

	//记录创建时的分数，默认为创建的时间
	pipe.ZAdd(ctx, KeyPostScoreZset, redis.Z{
		Score:  float64(post.CreateTime.Unix()),
		Member: fmt.Sprintf("%v", post.ID),
	})

	//记录在对应的community的set下面
	key := KeyPostCommunitySetPrefix + strconv.FormatInt(post.CommunityID,10)
	pipe.SAdd(ctx, key,fmt.Sprintf("%v",post.ID) )

	_, err := pipe.Exec(ctx)

	if err != nil {
		zap.L().Error("redis add KeyPostTimeZset or KeyPostScoreZset failed", zap.Error(err), zap.Int64("post_id", post.ID))
	}

	return nil
}
