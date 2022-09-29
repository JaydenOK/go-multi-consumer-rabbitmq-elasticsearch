package controllers

import (
	"app/services"
	"app/utils"
	"github.com/gin-gonic/gin"
)

type EsController struct {
	esService services.EsService
}

func (c *EsController) IndexLists(ctx *gin.Context) {
	returnData := c.esService.IndexLists(ctx)
	if returnData == nil {
		utils.FailResponse(ctx, nil)
		return
	}
	utils.SuccessResponse(ctx, returnData)
}

func (c *EsController) IndexExist(ctx *gin.Context) {
	returnData := c.esService.IndexExist(ctx)
	if returnData == nil {
		utils.FailResponse(ctx, nil)
		return
	}
	utils.SuccessResponse(ctx, returnData)
}

func (c *EsController) IndexCreate(ctx *gin.Context) {
	returnData := c.esService.IndexCreate(ctx)
	if returnData == nil {
		utils.FailResponse(ctx, nil)
		return
	}
	utils.SuccessResponse(ctx, returnData)
}

func (c *EsController) IndexGetMapping(ctx *gin.Context) {
	returnData := c.esService.IndexGetMapping(ctx)
	if returnData == nil {
		utils.FailResponse(ctx, nil)
		return
	}
	utils.SuccessResponse(ctx, returnData)
}

func (c *EsController) IndexPutMapping(ctx *gin.Context) {
	returnData := c.esService.IndexPutMapping(ctx)
	if returnData == nil {
		utils.FailResponse(ctx, nil)
		return
	}
	utils.SuccessResponse(ctx, returnData)
}

func (c *EsController) IndexReindex(ctx *gin.Context) {
	returnData := c.esService.IndexReindex(ctx)
	if returnData == nil {
		utils.FailResponse(ctx, nil)
		return
	}
	utils.SuccessResponse(ctx, returnData)
}

func (c *EsController) IndexDelete(ctx *gin.Context) {
	returnData := c.esService.IndexDelete(ctx)
	if returnData == nil {
		utils.FailResponse(ctx, nil)
		return
	}
	utils.SuccessResponse(ctx, returnData)
}

func (c *EsController) IndexAliasLists(ctx *gin.Context) {
	returnData := c.esService.IndexAliasLists(ctx)
	if returnData == nil {
		utils.FailResponse(ctx, nil)
		return
	}
	utils.SuccessResponse(ctx, returnData)
}

func (c *EsController) IndexAlias(ctx *gin.Context) {
	returnData := c.esService.IndexAlias(ctx)
	if returnData == nil {
		utils.FailResponse(ctx, nil)
		return
	}
	utils.SuccessResponse(ctx, returnData)
}
