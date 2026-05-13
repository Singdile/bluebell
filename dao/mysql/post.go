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
}

// InsertPost 插入Post到数据库
func InsertPost(post *models.Post) error {
	sqlstr := "INSERT INTO post (post_id,title,content,author_id,community_id,status,create_time) VALUES (:post_id,:title,:content,:author_id,:community_id,:status,:create_time)"

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
	sqlstr := "select post_id,title,content,author_id,community_id,status,create_time from post where post_id = ?"

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

// GetPostList 获取包含简略content内容的Postlist post的总数
// 获取分页数据，使用left join 欻性能作者名和社区名
func GetPostList(page, pagesize int64) ([]*models.PostListItem, int64, error) {
	// 计算offset
	offset := (page - 1) * pagesize

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
p.create_time AS create_time
FROM post p
LEFT JOIN user u ON p.author_id = u.user_id
LEFT JOIN community c ON p.community_id = c.community_id
ORDER BY p.create_time DESC
LIMIT ? OFFSET ?
`

	// 执行查询，映射结构
	var dblist []postListItemDB
	if err := db.Select(&dblist, sqlstr, pagesize, offset); err != nil {
		zap.L().Error("query post list failed", zap.Error(err))
		return nil, 0, err
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
		})
	}

	//查询记录的总数
	totalsql := "SELECT COUNT(*) FROM post"
	var total int64
	if err := db.Get(&total, totalsql); err != nil {
		zap.L().Error("query post total failed", zap.Error(err))
		return nil, 0, err
	}

	return list, total, nil
}

// 根据给定的post_id list 查询帖子数据
func GetPostListByIDs(post_ids []string, community_id int64) ([]*models.PostListItem, int64, error) {
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
    p.create_time AS create_time
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
		return nil, 0, err
	}

	query = db.Rebind(query)

	err = db.Select(&postlistitem, query, args...)
	if err != nil {
		zap.L().Error("sqlx.Select() failed", zap.Error(err))
		return nil, 0, err // ← 注意：这里必须返回 err，不能返回 nil
	}

	// 统计总数
	// community_id = 0 统计全站
	// community_id > 0 统计社区
	var total int64
	if community_id == 0 {
		totalsql := "SELECT COUNT(*) FROM post"
		if err := db.Get(&total, totalsql); err != nil {
			zap.L().Error("query post total failed", zap.Error(err))
			return nil, 0, err
		}

	} else if community_id > 0 {
		totalsql := "SELECT COUNT(*) FROM post WHERE community_id = ?"
		if err := db.Get(&total, totalsql,community_id); err != nil {
			zap.L().Error("query community post number failed", zap.Error(err))
			return nil, 0, err
		}
	}

	return postlistitem, total, nil
}
