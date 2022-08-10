package middleware

import (
	"app/pkg/response"
	"app/pkg/token"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// 使用方式：Header中增加 Authorization: Bearer <token>
const authorization = "Authorization"

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := ctx.GetHeader(authorization)
		index := strings.IndexByte(auth, ' ')
		if auth == "" || index < 0 {
			logrus.WithField(authorization, auth).Debug("Authorization is empty")
			response.Echo(ctx, nil, response.Msg("Token为空")) //
			return
		}

		accessToken := auth[index+1:]

		claims, err := token.ParseJWT(accessToken)
		if err != nil {
			if claims.IsExpired() {
				response.Echo(ctx, nil, response.Code(2001)) // "Token已过期"
			} else {
				logrus.WithField(authorization, auth).Warn(err.Error())
				response.Echo(ctx, nil, response.Msg("Token错误"))
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
