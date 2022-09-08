package routers

import (
	"app/controllers"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	orderGroup := r.Group("/order")
	orderGroup.GET("list", controllers.Lists)
}
