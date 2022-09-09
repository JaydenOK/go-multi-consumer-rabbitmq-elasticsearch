package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	StatusSuccess = 1
	StatusFail    = 2
)

func JsonResponse(ctx *gin.Context, httpStatus, status int, message string, data interface{}) {
	ctx.JSON(httpStatus, gin.H{
		"data":    data,
		"message": message,
		"status":  status,
	})
}

func SuccessResponse(ctx *gin.Context, data interface{}) {
	JsonResponse(ctx, http.StatusOK, StatusSuccess, "success", data)
}

func FailResponse(ctx *gin.Context, data interface{}) {
	JsonResponse(ctx, http.StatusOK, StatusFail, "fail", data)
}

func MessageResponse(ctx *gin.Context, message string) {
	JsonResponse(ctx, http.StatusOK, StatusSuccess, "success", message)
}
