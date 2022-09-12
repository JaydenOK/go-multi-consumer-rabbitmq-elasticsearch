package controllers

import (
	"app/services"
	"app/utils"
	"github.com/gin-gonic/gin"
)

// order控制器类
type UserController struct {
	userService services.UserService
}

func (c *UserController) Register(ctx *gin.Context) {
	returnData := c.userService.Register(ctx)
	if returnData == nil {
		utils.FailResponse(ctx, nil)
		return
	}
	utils.SuccessResponse(ctx, returnData)
}

func (c *UserController) SignIn(ctx *gin.Context) {
	returnData := c.userService.SignIn(ctx)
	if returnData == nil {
		utils.FailResponse(ctx, nil)
		return
	}
	utils.SuccessResponse(ctx, returnData)
}

func (c *UserController) SignOut(ctx *gin.Context) {
	returnData := c.userService.SignOut(ctx)
	if returnData == nil {
		utils.FailResponse(ctx, nil)
		return
	}
	utils.SuccessResponse(ctx, returnData)
}

func (c *UserController) List(ctx *gin.Context) {
	returnData := c.userService.List(ctx)
	if returnData == nil {
		utils.FailResponse(ctx, nil)
		return
	}
	utils.SuccessResponse(ctx, returnData)
}
