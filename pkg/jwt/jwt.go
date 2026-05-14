package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 过期时间
const TokenExpireDuration = time.Hour * 2
const (
	errUserID = "userID must be positive"
	errUserName = "username cannot be empty"
)

// salt
var mySecret = []byte("singdile")

// 定义自己的payload
type MyClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenToken 生成JWT
func GenToken(userID int64, username string) (string, error) {
	// 参数校验
	if userID <= 0 {
		return "",errors.New(errUserID)
	}

	if username == ""{
		return "",errors.New(errUserName)
	}

	//创建自定义的payload
	c := MyClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)),
			Issuer:    "bluebell",
		},
	}

	//创建签名对象
	tokenobj := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	//添加salt  签名返回
	return tokenobj.SignedString(mySecret)
}

// ParseToken 解析并验证
func ParseToken(tokenString string) (*MyClaims, error) {
	var mc = new(MyClaims)

	//解析 Token
	// 第三个参数是一个回调函数，用于add salt
	token, err := jwt.ParseWithClaims(tokenString, mc, func(t *jwt.Token) (any, error) { return mySecret, nil })

	// 处理解析产生的错误，比如过期，签名错误，格式错误
	if err != nil {
		return nil, err
	}

	//最终检验的有效性
	if token.Valid {
		return mc, nil
	}

	//转化为自己的payload
	if mc, ok := token.Claims.(*MyClaims); ok {
		return mc, nil
	}

	return nil, errors.New("invalid token")
}
