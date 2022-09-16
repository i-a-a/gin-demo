package middleware

import (
	"app/pkg/carrot"
	"app/pkg/token"
	"strings"

	"github.com/gin-gonic/gin"
)

// 使用方式：Header中增加 Authorization: Bearer <token> , 这是一种规范
const authorization = "Authorization"

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := ctx.GetHeader(authorization)
		index := strings.IndexByte(auth, ' ')
		if auth == "" || index < 0 {
			// carrot.New(ctx).Echo(nil, types.Msg("Authorization is empty"))
			return
		}

		accessToken := auth[index+1:]
		claims, err := token.ParseJWT(accessToken)
		if err != nil {
			if claims.IsExpired() {
				// token过期，需要前端判断code，执行刷新token动作
				carrot.New(ctx).Echo(nil, carrot.Msg("Authorization is expired"))
			} else {
				carrot.New(ctx).Echo(nil, carrot.Msg("Authorization is wrong"))
			}
			return
		}

		if claims.Uid == 0 {
			panic("uid is 0")
		}

		// 上下文带上用户ID
		ctx.Set("uid", claims.Uid)
	}
}
