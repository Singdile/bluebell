package redis

// redis key
// redis key 使用命令空间的方式，方便查询和拆分
// KeyPostTimeZset member post_id score post_time
// KeyPostScoreZset member post_id score votescore
// KeyPostVotedZsetPrefix+post_id member user_id score vote_direction
const (
	KeyPostTimeZset        = "bluebell:post:time"  // 帖子以及发帖时间
	KeyPostScoreZset       = "bluebell:post:score" //帖子以及投票分数
	KeyPostVotedZsetPrefix = "bluebell:post:voted:" //用户及投票类型
)
