package mysql

import (
	"bluebell/models"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// postListItemDB 用于数据库映射结构体
// 用于接收join 之后的查询结果 ，仅在本文件中使用
type postListItemDB struct {
	PostID         int64     `db:"post_id"`
	Title          string    `db:"title"`
	ContentPreview string    `db:"content_preview"` //截取摘要
	AuthorID       int64     `db:"author_id"`
	AuthorName     string    `db:"author_name"`
	CommunityID    int64     `db:"community_id"`
	CommunityName  string    `db:"community_name"`
	Status         int32     `db:"status"`
	CreateTime     time.Time `db:"create_time"`
	VoteP          int64     `db:"vote_p"`
	VoteN          int64     `db:"vote_n"`
}

// postDeteailDB 用于数据库映射结构体
// 用于详情查询
type postDeteailDB struct {
	PostID        int64     `db:"post_id"`
	Title         string    `db:"title"`
	Content       string    `db:"content"` //完整内容
	AuthorID      int64     `db:"author_id"`
	AuthorName    string    `db:"author_name"`
	CommunityID   int64     `db:"community_id"`
	CommunityName string    `db:"community_name"`
	Introduction  string    `db:"introduction"`
	Status        int32     `db:"status"`
	CreateTime    time.Time `db:"create_time"`
	VoteP         int64     `db:"vote_p"`
	VoteN         int64     `db:"vote_n"`
}

// InsertPost 插入Post到数据库
func InsertPost(post *models.Post) error {
	sqlstr := "INSERT INTO post (post_id,title,content,author_id,community_id,status,create_time,score) VALUES (:post_id,:title,:content,:author_id,:community_id,:status,:create_time,:score)"

	// 使用NamedExec, sqlx根据结构体的db tag来匹配
	_, err := db.NamedExec(sqlstr, post)
	if err != nil {
		zap.L().Error("insert post failed", zap.Error(err))
		return err
	}
	return nil
}

// GetPostByID 根据id查询post信息
func GetPostByID(post_id int64) (post *models.Post, err error) {
	sqlstr := "select post_id,title,content,author_id,community_id,status,create_time,vote_p,vote_n from post where post_id = ?"

	post = &models.Post{}
	/// 查询数据库操作失败
	if err = db.Get(post, sqlstr, post_id); err != nil {
		zap.L().Error("query post failed", zap.Error(err))
		return nil, err
	}

	//查询数cc据库操作成功
	return post, nil
}

// GetPostDetailByID 获取postdetail信息
func GetPostDetailByID(post_id int64) (postdetail *models.PostDetail, err error) {
	//sql语句
	sqlstr := `
    SELECT 
        p.post_id AS post_id,
        p.title AS title,
        p.content AS content,
        p.author_id AS author_id,
        p.community_id AS community_id,
        p.status AS status,
        p.create_time AS create_time,
        p.vote_p AS vote_p,
        p.vote_n AS vote_n,
        u.username AS author_name,
        c.community_name AS community_name,
        c.introduction AS introduction
    FROM post p
    LEFT JOIN user u ON p.author_id = u.user_id
    LEFT JOIN community c ON p.community_id = c.community_id
    WHERE p.post_id = ?
`
	//查询
	var dbData postDeteailDB               // ✅ 使用中间结构体
	err = db.Get(&dbData, sqlstr, post_id) // ✅ 映射到扁平结构体
	if err != nil {                        // 失败
		zap.L().Error("GetPostDetailByID(post_id int64) failed", zap.Error(err), zap.Int64("post_id", post_id))
		return nil, err
	}

	// 手动组装嵌套结构体
	postdetail = &models.PostDetail{
		AuthorName: dbData.AuthorName,
		Post: &models.Post{
			ID:          dbData.PostID,
			Title:       dbData.Title,
			Content:     dbData.Content,
			AuthorID:    dbData.AuthorID,
			CommunityID: dbData.CommunityID,
			Status:      dbData.Status,
			CreateTime:  dbData.CreateTime,
			VoteP:       dbData.VoteP,
			VoteN:       dbData.VoteN,
		},
		Community: &models.CommunityDetail{
			ID:           dbData.CommunityID,
			Name:         dbData.CommunityName,
			Introduction: dbData.Introduction, // ✅ 修复拼写后可用
		},
	}

	//成功
	return postdetail, err

}

// GetPostList 获取包含简略content内容的Postlist
func GetPostListByOrder(page, pagesize int64, order string) ([]*models.PostListItem, error) {
	// 计算offset
	offset := (page - 1) * pagesize

	// 确定排序字段
	orderby := "create_time"
	if order == "score" {
		orderby = "score"
	}

	//sql 语句查询，截取摘要
	sqlstr := `
SELECT
p.post_id as post_id,
p.title as title,
SUBSTRING(p.content,1,200) AS content_preview,
p.author_id AS author_id,
u.username AS author_name,
p.community_id AS community_id,
c.community_name AS community_name,
p.status AS status,
p.create_time AS create_time,
p.vote_p AS vote_p,
p.vote_n AS vote_n
FROM post p
LEFT JOIN user u ON p.author_id = u.user_id
LEFT JOIN community c ON p.community_id = c.community_id
ORDER BY p.` + orderby + ` DESC
LIMIT ? OFFSET ?
`

	// 执行查询，映射结构
	var dblist []postListItemDB
	if err := db.Select(&dblist, sqlstr, pagesize, offset); err != nil {
		zap.L().Error("query post list failed", zap.Error(err))
		return nil, err
	}

	// 组装成PostListItem ,对外暴露
	list := make([]*models.PostListItem, 0, len(dblist))
	for _, item := range dblist {
		list = append(list, &models.PostListItem{
			PostID:         item.PostID,
			Title:          item.Title,
			ContentPreview: item.ContentPreview,
			AuthorName:     item.AuthorName,
			AuthorID:       item.AuthorID,
			CommunityID:    item.CommunityID,
			CommunityName:  item.CommunityName,
			Status:         item.Status,
			CreateTime:     item.CreateTime,
			VoteP:          item.VoteP,
			VoteN:          item.VoteN,
		})
	}

	return list, nil
}

// 根据给定的post_id list 查询帖子数据
func GetPostListByIDs(post_ids []string, community_id int64) ([]*models.PostListItem, error) {
	idsStr := strings.Join(post_ids, ",")

	sqlstr := `
    SELECT
    p.post_id as post_id,
    p.title as title,
    SUBSTRING(p.content,1,200) AS content_preview,
    p.author_id AS author_id,
    u.username AS author_name,
    p.community_id AS community_id,
    c.community_name AS community_name,
    p.status AS status,
    p.create_time AS create_time,
    p.vote_p AS vote_p,
    p.vote_n AS vote_n
    FROM post p
    LEFT JOIN user u ON p.author_id = u.user_id
    LEFT JOIN community c ON p.community_id = c.community_id
    WHERE p.post_id IN (?)
    ORDER BY FIND_IN_SET(p.post_id, ?)
    `
	postlistitem := make([]*models.PostListItem, 0, len(post_ids))
	query, args, err := sqlx.In(sqlstr, post_ids, idsStr)
	if err != nil {
		zap.L().Error("sqlx.In() failed", zap.Error(err))
		return nil, err
	}

	query = db.Rebind(query)

	err = db.Select(&postlistitem, query, args...)
	if err != nil {
		zap.L().Error("sqlx.Select() failed", zap.Error(err))
		return nil, err // ← 注意：这里必须返回 err，不能返回 nil
	}

	return postlistitem, nil
}

// GetTotalPostCount 获取帖子的总数
func GetTotalPostCount() (int64, error) {
	var total int64
	err := db.Get(&total, "SELECT COUNT(*) FROM post")
	if err != nil {
		zap.L().Error("Failed to get total post count from mysql", zap.Error(err))
		return 0, err
	}

	return total, nil
}

// GetCommunityPostCount 获取社区帖子总数
func GetCommunityPostCount(communityid int64) (int64, error) {
	var total int64
	err := db.Get(&total, "SELECT COUNT(*) FROM post WHERE community_id= ?", communityid)
	if err != nil {
		zap.L().Error("Failed to get community post count from mysql", zap.Error(err), zap.Int64("community_id", communityid))
		return 0, err
	}

	return total, nil

}

// postForSync 用于同步 Redis 索引的帖子数据
// 只包含索引需要的字段
type postForSync struct {
	PostID      int64 `db:"post_id"`
	CommunityID int64 `db:"community_id"`
	CreateTime  int64 `db:"create_time"` // 时间戳
	Score       int64 `db:"score"`
}

// PostForSync 导出的帖子同步结构体
type PostForSync struct {
	PostID      int64
	CommunityID int64
	CreateTime  int64
	Score       int64
}

// GetAllPostsForSync 获取所有帖子信息，用于同步 Redis 索引
// 只查询索引需要的字段：post_id, community_id, create_time, score
func GetAllPostsForSync() ([]*PostForSync, error) {
	sqlstr := `SELECT post_id, community_id, UNIX_TIMESTAMP(create_time) AS create_time, score FROM post`

	var posts []postForSync
	if err := db.Select(&posts, sqlstr); err != nil {
		zap.L().Error("Failed to get all posts for sync",
			zap.Error(err))
		return nil, err
	}

	// 转换为导出结构体
	result := make([]*PostForSync, 0, len(posts))
	for _, p := range posts {
		result = append(result, &PostForSync{
			PostID:      p.PostID,
			CommunityID: p.CommunityID,
			CreateTime:  p.CreateTime,
			Score:       p.Score,
		})
	}

	return result, nil
}
