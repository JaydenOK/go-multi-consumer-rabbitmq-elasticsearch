package controllers

import (
	"app/lib/mysql"
	"app/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func Lists(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	orderId := c.Query("order_id")
	platformCode := c.Query("platform_code")
	if page < 0 {
		page = 1
	}
	if pageSize < 0 {
		pageSize = 50
	}
	var order models.Order		//用于查找单个
	var orders []models.Order		//用于查找多个
	db := mysql.GetMysqlClient().Table(order.TableName())
	if orderId != "" {
		db = db.Where("order_id=?", orderId)
	}
	if platformCode != "" {
		platformCodeList := strings.Split(platformCode, ",")
		db = db.Where("platform_code IN ?", platformCodeList)
	}
	db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&orders)
	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "ok",
		"data":    orders,
	})
}
