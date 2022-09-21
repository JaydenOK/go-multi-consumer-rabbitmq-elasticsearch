package controllers

import (
	"app/services"
	"app/utils"
	"github.com/gin-gonic/gin"
)

// order控制器类
type OrderController struct {
	orderService services.OrderService
}

func (c *OrderController) Lists(ctx *gin.Context) {
	returnData := c.orderService.Lists(ctx)
	if returnData == nil {
		utils.FailResponse(ctx, nil)
		return
	}
	utils.SuccessResponse(ctx, returnData)
}

// es订单列表
func (c *OrderController) EsLists(ctx *gin.Context) {
	returnData := c.orderService.EsLists(ctx)
	if returnData == nil {
		utils.FailResponse(ctx, nil)
		return
	}
	utils.SuccessResponse(ctx, returnData)
}

func (c *OrderController) Add(ctx *gin.Context) {
	returnData := c.orderService.Add(ctx)
	if returnData == nil {
		utils.FailResponse(ctx, nil)
		return
	}
	utils.SuccessResponse(ctx, returnData)
}

func (c *OrderController) Update(ctx *gin.Context) {
	returnData := c.orderService.Update(ctx)
	if returnData == nil {
		utils.FailResponse(ctx, nil)
		return
	}
	utils.SuccessResponse(ctx, returnData)
}

func (c *OrderController) Delete(ctx *gin.Context) {
	returnData := c.orderService.Delete(ctx)
	if returnData == nil {
		utils.FailResponse(ctx, nil)
		return
	}
	utils.SuccessResponse(ctx, returnData)
}
