package redis

import (
	"bluebell/dao/mysql"
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// SyncPostsFromMySQL 批量同步mysql帖子索引到 reddis
func SyncPostsFromMySQL(ctx context.Context) error {
	// 从mysql中获取帖子索引信息
	posts, err := mysql.GetAllPostsForSync()
	if err != nil {
		zap.L().Error("Failed to get posts from MySQL for sync", zap.Error(err))
		return err
	}

	if len(posts) == 0 { //mysql里面也没有数据
		return nil
	}

	zap.L().Info("Syncing Redis index from MySQL")
	// redis 同步开始
	// 使用pipeline 批量写入
	pipeline := rdb.TxPipeline()

	for _, post := range posts {
		postidstr := strconv.FormatInt(post.PostID, 10)
		communityKey := KeyPostCommunitySetPrefix + strconv.FormatInt(post.CommunityID,10)
		//时间索引ZSET
		pipeline.ZAdd(ctx, KeyPostTimeZset, redis.Z{
			Score:  float64(post.CreateTime),
			Member: postidstr,
		})
		//投票分数索引ZSET
		pipeline.ZAdd(ctx, KeyPostScoreZset, redis.Z{
			Score:  float64(post.Score),
			Member: postidstr,
		})

		//社区索引SET
		pipeline.SAdd(ctx, communityKey, postidstr)

	}

	_, err = pipeline.Exec(ctx)
	if err != nil {
		zap.L().Error("Failed to sync posts to Redis", zap.Error(err))
		return err
	}
	return nil

}
