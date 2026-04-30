package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// 过期时间
const TokenExpireDuration = time.Hour * 2

// salt
var mySecret = []byte("singdile")

type MyClaims struct {
	UserID int64 `json:"user_id"`
	jwt.StandardClaims
}
