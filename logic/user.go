package logic

import (
	"bluebell/dao/mysql"
)

func SignUp() (err error) {
	//查询用户是否已经存在

	//插入用户数据
	mysql.InsertUser()
	return nil

}
