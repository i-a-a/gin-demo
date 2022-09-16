package middleware

import (
	"app/pkg/carrot"
	"fmt"

	"github.com/gin-gonic/gin"
)

func CustomRecovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		err, ok := recovered.(error)
		if !ok {
			err = fmt.Errorf("%v", recovered)
		}
		carrot.New(c).Echo(nil, err)
	})
}
