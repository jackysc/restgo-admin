package restgo

import (
	"github.com/gin-gonic/gin"
)

func IsLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := LoadUser(ctx)
		if user == nil {
			ResultFail(ctx, "尚未登录")
			ctx.Abort()
		}
	}
}
