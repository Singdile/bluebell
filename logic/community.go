package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
)

// GetCommunityList
func GetCommunityList() ([]*models.Community, error) {
	// 查询所有的community数据
	data, err := mysql.GetCommunityList()

	return data, err

}

// GetCommunityByID 根据id获取社区详情
func GetCommunityByID(id int64) (communitydetail *models.CommunityDetail, err error) {
	communitydetail, err = mysql.GetCommunityByID(id)
	return
}
