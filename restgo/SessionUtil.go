package restgo

import (
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SetSession(ctx *gin.Context, k string, o interface{}) {
	session := sessions.Default(ctx)
	session.Set(k, o)
	err := session.Save()
	if err != nil {
		fmt.Println(err)
	}
}

func GetSession(ctx *gin.Context, k string) interface{} {
	session := sessions.Default(ctx)
	return session.Get(k)
}

func SaveUser(ctx *gin.Context, user interface{}) {
	SetSession(ctx, "user", user)
}

func LoadUser(ctx *gin.Context) interface{} {
	return GetSession(ctx, "user")
}

func SaveRoleId(ctx *gin.Context, roleId interface{}) {
	session := sessions.Default(ctx)

	session.Set("roleid", roleId)
	session.Save()
}

func LoadRoleId(ctx *gin.Context) interface{} {
	session := sessions.Default(ctx)
	o := session.Get("roleid")
	return o

}

func ClearAllSession(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	session.Save()
	return
}
