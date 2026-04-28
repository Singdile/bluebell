package mysql

import (
	"bluebell/models"

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

// 查询用户
func SelectUserByName() {}

// 查询用户是否存在
func CheckUserExist(username string) (bool, error) {
	sqlstr := "select count(user_id) from user where username = ?"
	var count int
	if err := db.Get(&count, sqlstr, username); err != nil {
		return false, err
	}

	return count > 0, nil
}
