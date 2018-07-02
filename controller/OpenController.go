package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jackysc/restgo-admin/restgo"
)

type OpenController struct {
	restgo.Controller
}

func (ctrl *OpenController) Router(router *gin.Engine) {

	r := router.Group("open")
	r.GET("verify", ctrl.verify)

}

func (ctrl *OpenController) verify(ctx *gin.Context) {
	restgo.LoadVerify(ctx)
}
