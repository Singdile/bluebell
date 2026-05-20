package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"errors"
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

// CreateCommunity 创建社区
func CreateCommunity(c *models.Community) (int64, error) {
	// 检查社区名称是否已经存在
	exist,err := mysql.CheckCommunityName(c.Name)
	if err != nil { //数据库操作有误
		return 0, err
	}
	if exist {
		return 0,errors.New("社区名称已经存在")
	}

	// 创建社区
	var id int64
	if id, err = mysql.CreateCommunity(c); err != nil {
		return 0, err
	}

	return id, nil
}

// UpdateCommunity 更新社区的信息
func UpdateCommunity(communityInfo *models.Community)  error {
	// 判断 id ，name 是否合法
	err := mysql.CheckCommunityIDExist(communityInfo.ID)
	if err != nil { //sql执行出现问题/或者id 不存在
		return err
	}

	exist,err := mysql.CheckCommunityNameExcludeID(communityInfo.Name,communityInfo.ID)
	if err != nil { //数据库错误
		return err
	}

	if exist { //除了自己以外该名称存在，则无法更新
		return errors.New("该名称存在，无法更新")
	}

	// 执行更新
	err = mysql.UpdateCommunity(communityInfo)
	if err != nil {
		return  err
	}

	return nil
}
