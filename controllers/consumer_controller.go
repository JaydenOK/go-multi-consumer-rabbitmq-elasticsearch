package controllers

import (
	"app/services"
	"app/utils"
	"github.com/gin-gonic/gin"
)

type ConsumerController struct {
	consumerService services.ConsumerService
}

// 接收
func (c *ConsumerController) StartConsumer(ctx *gin.Context) {
	returnData := c.consumerService.StartConsumer(ctx)
	if returnData == nil {
		utils.FailResponse(ctx, nil)
		return
	}
	utils.SuccessResponse(ctx, returnData)
}

func (c *ConsumerController) StopConsumer(ctx *gin.Context) {
	returnData := c.consumerService.StopConsumer(ctx)
	if returnData == nil {
		utils.FailResponse(ctx, nil)
		return
	}
	utils.SuccessResponse(ctx, returnData)
}

func (c *ConsumerController) StopAll(ctx *gin.Context) {
	returnData := c.consumerService.StopAll(ctx)
	if returnData == nil {
		utils.FailResponse(ctx, nil)
		return
	}
	utils.SuccessResponse(ctx, returnData)
}
