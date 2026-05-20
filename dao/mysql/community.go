package mysql

import (
	"bluebell/models"
	"database/sql"
	"errors"

	"go.uber.org/zap"
)

// GetCommunityList 查询数据库中的社区数据，返回列表数据
func GetCommunityList() (data []*models.Community, err error) {
	sqlstr := "select id, community_name from community"

	err = db.Select(&data, sqlstr)

	if err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no community in db")
			err = nil
		}
	}
	return
}

// GetCommunityByID 根据id获取社区的详情
func GetCommunityByID(id int64) (communitydetail *models.CommunityDetail, err error) {
	communitydetail = new(models.CommunityDetail)

	sqlstr := "select id, community_name, introduction from community where id = ?"

	err = db.Get(communitydetail, sqlstr, id)

	if err != nil {
		if err == sql.ErrNoRows { //特殊情况，id不存在，记录日志
			zap.L().Warn("communit not found", zap.Int64("communit_id", id))
			return nil, nil
		}
		// 一般情况，查询执行失败
		zap.L().Error("query community by id failed", zap.Int64("community_id", id))

		return nil, err
	}

	// 查询成功
	return communitydetail, nil
}

// CheckCommunityName 是否已经存在
// bool -true 存在， false - 不存在 ， error - 数据库错误
func CheckCommunityName(name string) (bool, error) {
	sqlstr := "select 1 from community where community_name = ? limit 1"
	var count int
	err := db.Get(&count,sqlstr,name)

	if err == sql.ErrNoRows {//表明查询成功，名称不存在
		return false,nil
	}

	if err != nil {//查询失败，数据库错误
		zap.L().Error("Failed to check community name", zap.Error(err))
		return false,err
	}

	return true, nil //查询成功

}

// CheckCommunityNameExcludeID 检查除了自己以外的社区名称是否重复
// true-存在，false-不存在
// error-数据库错误
func CheckCommunityNameExcludeID(name string, excludeID int64) (bool, error) {
	sqlstr := "SELECT 1 FROM community WHERE community_name = ? AND id != ? LIMIT 1"
	var exists int
	err := db.Get(&exists, sqlstr, name, excludeID)

	if err == sql.ErrNoRows { //查询到名称不存在
		return false, nil
	}

	if err != nil { //数据库操作错误
		zap.L().Error("Failed to check community name", zap.Error(err))
		return false,err
	}

	return exists==1,nil //查询到名称，存在
}

// CreateCommunity 创建社区
func CreateCommunity(c *models.Community) (int64, error) {
	// sql语句, 采用sqlx的命名查询，这里: 后面实际是对应的结构体的db标签
	sqlstr := "INSERT INTO community (community_name,introduction,creator_id,status) VALUES (:community_name,:introduction,:creator_id,:status)"
	// 执行sql
	res, err := db.NamedExec(sqlstr, c)
	if err != nil {
		zap.L().Error("create community in mysql failed", zap.Error(err), zap.String("name", c.Name))
		return 0, err
	}
	// 返回
	lastinsertid, err := res.LastInsertId()
	if err != nil {
		zap.L().Error("Failed to get last insert community id", zap.Error(err))
		return 0, err
	}
	return lastinsertid, nil
}

// CheckIDExist 检查id是否存在
// TODO: 优化风格一致 返回类型 (bool,error)
func CheckCommunityIDExist(id int64) error {
	//sql
	sqlstr := "select 1 from community where id = ? limit 1"

	// 执行
	var count int
	err := db.Get(&count, sqlstr, id)
	if err != nil {
		zap.L().Error("mysql failed to query community id exists or not ", zap.Error(err), zap.Int64("community_id", id))
		return errors.New("查询community id 是否存在出现错误")
	}

	if count != 1 {
		return errors.New("community id 不存在")
	}

	return nil
}

// UpdateCommunity 更新社区数据
func UpdateCommunity(community *models.Community) error {
	// sql
	sqlstr := "UPDATE community SET community_name = ?,introduction = ? WHERE id = ?"

	// 执行
	_, err := db.Exec(sqlstr, community.Name, community.Introduction, community.ID)
	if err != nil {
		zap.L().Error("update community failed", zap.Error(err), zap.Int64("id", community.ID))
		return err
	}

	return nil
}
