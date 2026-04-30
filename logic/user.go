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

// Login 检查用户是否存在，密码是否正确，是否能够登录
func Login(p *models.ParamLogin) (err error) {
	//检查用户是否存在
	exist, err := mysql.CheckUserExist(p.Username)

	if !exist {
		return errors.New("用户不存在")
	}

	//验证密码是否一致
	if err = mysql.VerifyUserLogin(p.Username, p.Password); err != nil {
		return errors.New("验证失败")
	} else {
		return nil
	}
}
