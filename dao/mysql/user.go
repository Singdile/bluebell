package mysql

import (
	"bluebell/models"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// InsertUser 插入新的用户记录
func InsertUser(user *models.User) (err error) {
	//密码hash保存
	user.Password, err = encryptionPassword(user.Password)
	if err != nil {
		return err
	}

	//执行sql语句插入用户数据
	sqlstr := "insert into user (user_id,username,password) values (?,?,?)"
	_, err = db.Exec(sqlstr, user.UserID, user.Username, user.Password)

	if err != nil {
		return err
	}

	return nil
}

// GetUserByName 根据用户名获取用户信息,user_id,username,password
func GetUserByName(username string) (*models.User, error) {
	sqlstr := "select user_id,username,password from user where username = ?"

	var user = new(models.User)
	if err := db.Get(user, sqlstr, username); err != nil {
		return nil, err
	}

	return user, nil
}

// encryptionPassword 计算得到hash之后的密码
func encryptionPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// checkPasswordHash 比较密码和数据库中的hash之后的密码
func checkPasswordHash(password, hash string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false
	}
	return true
}

// VerifyUserLogin 验证用户登录信息与数据库记录是否一致
func VerifyUserLogin(param *models.ParamLogin) (err error) {
	//查询用户完整信息
	user, err := GetUserByName(param.Username)
	if err != nil {
		return errors.New("用户不存在")
	}

	//验证密码是否一致
	if checkPasswordHash(param.Password, user.Password) {
		//验证通过，修改param ，返回对应的user_id
		param.UserID = user.UserID
		return nil
	} else {
		return errors.New("密码错误,验证失败")
	}
}

// 查询用户
func getPasswordByName(username string) (string, error) {
	sqlstr := "select password from user where username = ?"
	password := new(string)

	if err := db.Get(&password, sqlstr, username); err != nil {
		return "", errors.New("failed to get password")
	}

	return *password, nil
}

// 查询用户是否存在
func CheckUserExist(username string) (bool, error) {
	sqlstr := "select count(user_id) from user where username = ?"
	var count int
	if err := db.Get(&count, sqlstr, username); err != nil {
		return false, err
	}

	return count > 0, nil
}
