package middlewares

import (
	"bluebell/controllers"
	"bluebell/pkg/jwt"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	jwt5 "github.com/golang-jwt/jwt/v5"
)

func JwtAuthMiddleware() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		//获取请求头的验证token
		authhead := ctx.Request.Header.Get("Authorization")
		if authhead == "" {
			//提示未认证
			controllers.FailWithDefault(ctx, controllers.ErrUnauthorized)
			ctx.Abort()
			return
		}

		//切割获取 tokenstring
		authstring := strings.SplitN(authhead, " ", 2)
		if !(len(authstring) == 2 && authstring[0] == "Bearer") {
			controllers.FailWithDefault(ctx, controllers.ErrUnauthorized)
			ctx.Abort()
			return
		}
		var claim = new(jwt.MyClaims)

		//解析验证tonken
		claim, err := jwt.ParseToken(authstring[1])

		if err != nil { //解析失败
			if errors.Is(err, jwt5.ErrTokenExpired) { //token过期
				controllers.FailWithDefault(ctx, controllers.ErrTokenExpired)
			} else { //token无效
				controllers.FailWithDefault(ctx, controllers.ErrInvalidToken)
			}
			ctx.Abort()
			return
		}

		//解析成功
		ctx.Set(controllers.Ctxuserid, claim.UserID) //保存userid
		ctx.Next()                                   // 后续可以通过ctxuserid,来获取当前请求的用户信息
	}

}
