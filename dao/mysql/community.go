package mysql

import (
	"bluebell/models"
	"database/sql"

	"go.uber.org/zap"
)

// GetCommunityList 查询数据库中的社区数据，返回列表数据
func GetCommunityList() (data []*models.Community, err error) {
	sqlstr := "select community_id, community_name from community"

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

	sqlstr := "select community_id, community_name, introduction from community where community_id = ?"

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
