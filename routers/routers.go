package routers

import (
	"app/controllers"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	//区块，无特殊意义
	orderGroup := r.Group("/order")
	{
		var orderController controllers.OrderController
		orderGroup.GET("list", orderController.Lists)
		orderGroup.GET("esList", orderController.EsLists)
		orderGroup.POST("add", orderController.Add)
		orderGroup.POST("update", orderController.Update)
		orderGroup.POST("delete", orderController.Delete)
	}

	//用户相关
	userGroup := r.Group("/user")
	{
		var userController controllers.UserController
		userGroup.POST("register", userController.Register)
		userGroup.GET("list", userController.List)
		userGroup.POST("signIn", userController.SignIn)
		userGroup.POST("signOut", userController.SignOut)
	}

}
