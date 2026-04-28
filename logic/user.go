package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/snowflake"
	"errors"
)

func SignUp(p *models.ParamSignUp) (err error) {
	//查询用户是否已经存在
	exist, err := mysql.CheckUserExist(p.Username)

	if exist {
		return errors.New("用户已经存在")
	}

	if err != nil {
		//数据库查询错误
		return err
	}

	//插入用户数据
	//1.构造一个user实例
	//  生成userid
	userid := snowflake.GenID()
	u := models.User{
		UserID:   userid,
		Username: p.Username,
		Password: p.Password,
	}
	//2.保存进入数据库
	if err := mysql.InsertUser(&u); err != nil {
		return err
	}
	return nil
}
