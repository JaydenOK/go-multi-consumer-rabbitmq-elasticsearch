package services

import (
	"app/lib/mysqllib"
	"app/models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

type OrderService struct {
}

// 订单查询
func (orderService *OrderService) Lists(ctx *gin.Context) interface{} {
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("pageSize"))
	orderId := ctx.Query("order_id")
	platformCode := ctx.Query("platform_code")
	middleCreateTimeStart := ctx.Query("middle_create_time_start")
	middleCreateTimeEnd := ctx.Query("middle_create_time_end")
	if page < 0 {
		page = 1
	}
	if pageSize < 0 {
		pageSize = 50
	}
	var order models.OrderModel    //用于查找单个
	var orders []models.OrderModel //用于查找多个
	db := mysqllib.GetMysqlClient().Table(order.TableName())
	if orderId != "" {
		db = db.Where("order_id=?", orderId)
	}
	//平台code查询，支持逗号分隔
	if platformCode != "" {
		platformCodeList := strings.Split(platformCode, ",")
		db = db.Where("platform_code IN ?", platformCodeList)
	}
	if middleCreateTimeStart != "" {
		db = db.Where("middle_create_time > ?", middleCreateTimeStart)
	}
	if middleCreateTimeEnd != "" {
		db = db.Where("middle_create_time <= ?", middleCreateTimeEnd)
	}
	db.Offset((page - 1) * pageSize).Limit(pageSize).Order("id desc").Find(&orders)
	return orders
}

// 新增订单 orderModel指定的属性
func (orderService *OrderService) Add(ctx *gin.Context) interface{} {
	var orderModel models.OrderModel
	if err := ctx.ShouldBind(&orderModel); err != nil {
		fmt.Println("bind error", orderModel)
		return nil
	}
	mysqlClient := mysqllib.GetMysqlClient()
	result := mysqlClient.Create(&orderModel) // 通过数据的指针来创建
	if result.Error != nil {
		fmt.Println(result.Error)
		return "新增订单失败:" + result.Error.Error()
	}
	return "新增订单成功，id为：" + strconv.Itoa(int(orderModel.Id))
}

// 通过order_id更新订单信息
func (orderService *OrderService) Update(ctx *gin.Context) interface{} {
	var orderModel models.OrderModel
	byteData, _ := ctx.GetRawData()
	if err := json.Unmarshal(byteData, &orderModel); err != nil {
		return "数据解析异常，请核对：" + err.Error()
	}
	obj := make(map[string]interface{})
	if err := json.Unmarshal(byteData, &obj); err != nil {
		return "数据解析异常，请核对：" + err.Error()
	}
	fmt.Println(obj)
	mysqlClient := mysqllib.GetMysqlClient()
	//批量更新
	result := mysqlClient.Model(&models.OrderModel{}).Where("order_id = ?", obj["order_id"]).Updates(obj)
	return "更新订单成功，id为：" + strconv.Itoa(int(result.RowsAffected))
}

// 通过order_id删除订单
func (orderService *OrderService) Delete(ctx *gin.Context) interface{} {
	var orderModel models.OrderModel
	orderId := ctx.PostForm("order_id")
	if orderId == "" {
		return ""
	}
	mysqlClient := mysqllib.GetMysqlClient()
	result := mysqlClient.First(&orderModel, "order_id = ?", orderId)
	if result.Error != nil || result.RowsAffected == 0 {
		fmt.Println(result.Error)
		return "订单不存在:" + result.Error.Error()
	}
	//批量删除
	result = mysqlClient.Where("order_id = ?", orderId).Delete(&models.OrderModel{})
	return "删除订单成功:" + strconv.Itoa(int(result.RowsAffected))
}
