package routes

import (
	"github.com/warlockz/ase-service/controller"

	"github.com/gin-gonic/gin"
)

func RouteIndex(Router *gin.Engine) {
	Router.POST("/encryptFile", controller.EncryptFile)
	Router.POST("/getFile", controller.GetFile)
}
