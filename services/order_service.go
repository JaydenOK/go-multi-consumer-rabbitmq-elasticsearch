package services

import (
	"app/constants"
	"app/events"
	"app/libs/elasticsearchlib"
	"app/libs/mysqllib"
	"app/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"strings"
)

type OrderService struct {
	eventManager events.EventManager
}

// 订单查询
func (orderService *OrderService) Lists(ctx *gin.Context) interface{} {
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("pageSize"))
	orderId := ctx.Query("order_id")
	platformCode := ctx.Query("platform_code")
	middleCreateTimeStart := ctx.Query("middle_create_time_start")
	middleCreateTimeEnd := ctx.Query("middle_create_time_end")
	totalPriceStart := ctx.Query("total_price_start")
	totalPriceEnd := ctx.Query("total_price_end")
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
	if totalPriceStart != "" {
		db = db.Where("total_price > ?", totalPriceStart)
	}
	if totalPriceEnd != "" {
		db = db.Where("total_price <= ?", totalPriceEnd)
	}
	db.Offset((page - 1) * pageSize).Limit(pageSize).Order("id desc").Find(&orders)
	return orders
}

// es订单列表
func (orderService *OrderService) EsLists(ctx *gin.Context) interface{} {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "0"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	orderId := ctx.Query("order_id")
	shipName := ctx.Query("ship_name")
	orderStatus := ctx.Query("order_status")
	platformCode := ctx.Query("platform_code")
	totalPriceStart := ctx.Query("total_price_start")
	totalPriceEnd := ctx.Query("total_price_end")
	middleCreateTimeStart := ctx.Query("middle_create_time_start")
	middleCreateTimeEnd := ctx.Query("middle_create_time_end")

	//fmt.Printf("pageSize: %+v", pageSize)
	//fmt.Printf("orderId: %+v", orderId)
	//fmt.Printf("shipName: %+v", shipName)
	//fmt.Printf("orderStatus: %+v", orderStatus)
	//fmt.Printf("platformCode: %+v", platformCode)
	//fmt.Printf("totalPriceStart: %+v,totalPriceEnd:%+v", totalPriceStart, totalPriceEnd)
	//fmt.Printf("middleCreateTimeStart: %+v,middleCreateTimeEnd:%+v", middleCreateTimeStart, middleCreateTimeEnd)

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 50
	}

	esClient := elasticsearchlib.GetClient()

	var res *esapi.Response
	var err error
	var r map[string]interface{}
	index := "order" //查询Index库
	//构造请求参数体
	//query参数
	//多条件，且，{"query":{"bool":{"must":[{"match":{"title":"test6"}},{"match":{"num":5}}]}}}
	//多条件，或，{"query":{"bool":{"should":[{"match":{"title":"test6"}},{"match":{"title":"test8"}}]}}}
	//查询 “num”小于4的数据 gt：大于 gte：大于等于 lt：小于 lte：小于等于,{"query":{"bool":{"filter":{"range":{"num":{"lt":4}}}}}}
	//范围，{"range":{"birthday":{"from":"1990-10-10","to":"2000-05-01","include_lower":true,"include_upper":false}}}
	//数值范围：{"query":{"range":{"price":{"gte":40,"lte":80,"boost":2}}}}
	//日期范围：{"query":{"range":{"post_date":{"gte":"2018-01-01 00:00:00","lte":"now","format":"yyyy-MM-dd hh:mm:ss","time_zone":"+1:00"}}}}

	//term 精确查询：只匹配指定的字段中包含指定的词的文档，单个。另外terms可指定多个字段	{"terms":{"price":[20,30]}}
	//term 是代表完全匹配，也就是精确查询，搜索前不会再对搜索词进行分词，所以我们的搜索词必须是文档分词集合中的一个	{"query":{"term":{"content":"中国"}}}
	//match 查询会先对搜索词进行分词,分词完毕后再逐个对分词结果进行匹配，因此相比于term的精确搜索，match是分词匹配搜索，match 查询相当于模糊匹配,只包含其中一部分关键词就行 {"match":{"a":123}}
	//match_all 能够匹配索引中的所有文件。 可以在查询中使用boost包含加权值，它将赋给所有跟它匹配的文档,计算score时用到。 {"match_all":{"a":123}}

	requestData := make(map[string]interface{})
	query := make(map[string]interface{})
	boolean := make(map[string]interface{}) //且查找
	//sortList := make([]interface{}, 5)      //排序
	//must := make([]interface{}, 20)         //多条件
	var sortList []interface{}
	var must []interface{}   //且查找
	var should []interface{} //或查找

	//分页
	requestData["from"] = (page - 1) * pageSize
	requestData["size"] = pageSize
	//排序: {"id":{"order":"desc"}}  （其中一个）
	idSort := map[string]map[string]string{
		"id": {
			"order": "desc",
		},
	}
	sortList = append(sortList, idSort)
	requestData["sort"] = sortList

	//精确查询
	if orderId != "" {
		//{"term":{"order_id":"orderId"}}
		orderSearch := map[string]interface{}{
			"term": map[string]string{
				"order_id": orderId,
			},
		}
		must = append(must, orderSearch)
	}

	//模糊分词查找
	if shipName != "" {
		//{"match":{"ship_name":"shipName"}}
		shipNameSearch := map[string]interface{}{
			"match": map[string]string{
				"ship_name": shipName,
			},
		}
		must = append(must, shipNameSearch)
	}

	//精确查询
	if orderStatus != "" {
		orderStatusSearch := map[string]interface{}{
			"term": map[string]string{
				"order_status": orderStatus,
			},
		}
		must = append(must, orderStatusSearch)
	}

	//多条件等值查询，平台code查询，支持逗号分隔，terms 可指定多个字段
	if platformCode != "" {
		platformCodeList := strings.Split(platformCode, ",")
		//查不到？？
		platformCodeSearch := map[string]interface{}{
			"terms": map[string]interface{}{
				"platform_code": platformCodeList,
			},
		}
		must = append(must, platformCodeSearch)
		//
		//使用term形式在es中query或者aggs时，发现竟然查不到数据，反而match却可以。几经波折，发现，如果字段type是text类型，是不支持term形式的keyword方式查找的，
		// 重新mapping，可以在该字段下加入keyword的fields，就可以实现term查找了。
		//platformCodeSearch := map[string]interface{}{
		//	"terms": map[string]interface{}{
		//		"platform_code": platformCodeList,
		//	},
		//}
		//should = append(should, platformCodeSearch)

	}

	//数值区间 {"query":{"range":{"price":{"gte":40,"lte":80,"boost":2}}}}
	if totalPriceStart != "" && totalPriceEnd != "" {
		totalPriceRange := map[string]interface{}{
			"range": map[string]interface{}{
				"total_price": map[string]interface{}{
					"gte": totalPriceStart,
					"lte": totalPriceEnd,
				},
			},
		}
		must = append(must, totalPriceRange)
	}

	//时间区间
	if middleCreateTimeStart != "" && middleCreateTimeEnd != "" {
		//{"range":{"birthday":{"from":"2022-09-28 13:51:27","to":"2022-12-11 04:42:00","include_lower":true,"include_upper":false}}}
		createTimeSearch := map[string]interface{}{
			"range": map[string]interface{}{
				"middle_create_time": map[string]interface{}{
					"from":          middleCreateTimeStart,
					"to":            middleCreateTimeEnd,
					"include_lower": true,
					"include_upper": false,
				},
			},
		}
		must = append(must, createTimeSearch)
	}

	boolean["must"] = must
	boolean["should"] = should
	query["bool"] = boolean
	requestData["query"] = query

	var testRequestData []byte
	if testRequestData, err = json.Marshal(requestData); err != nil {
		fmt.Println("error requestData:" + err.Error())
	}
	fmt.Println("ES请求数据：", string(testRequestData))
	////解析数据到buf，再请求es查询
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(requestData); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}
	res, err = esClient.Search(
		esClient.Search.WithContext(context.Background()),
		esClient.Search.WithIndex(index),
		esClient.Search.WithBody(&buf),
		//esClient.Search.WithBody(strings.NewReader(request)),
		esClient.Search.WithTrackTotalHits(true),
		esClient.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			log.Fatalf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}
	fmt.Println("4", err)
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// Print the response status, number of results, and request duration.
	log.Printf(
		"[%s] %d hits; took: %dms",
		res.Status(),
		int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		int(r["took"].(float64)),
	)
	var orders []interface{}
	dataList := map[string]interface{}{}
	dataList["page"] = page
	dataList["page_size"] = pageSize
	dataList["total"] = r["hits"].(map[string]interface{})["total"]
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		//log.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
		orders = append(orders, hit.(map[string]interface{})["_source"])
	}
	log.Println(strings.Repeat("=", 37))
	dataList["records"] = orders
	return dataList
}

// 新增订单 orderModel指定的属性
func (orderService *OrderService) Add(ctx *gin.Context) interface{} {
	orderService.registerEvent()
	var orderModel models.OrderModel
	if err := ctx.ShouldBind(&orderModel); err != nil {
		fmt.Println("bind error", err.Error())
		return "参数提交错误，请核对：" + err.Error()
	}
	mysqlClient := mysqllib.GetMysqlClient()
	result := mysqlClient.Create(&orderModel) // 通过数据的指针来创建
	if result.Error != nil {
		fmt.Println(result.Error)
		return "新增订单失败:" + result.Error.Error()
	}
	//发送mq通知程序，更新es信息
	orderService.eventManager.Trigger(constants.EventOrderChange, orderModel.OrderId)
	return "新增订单成功，id为：" + strconv.Itoa(int(orderModel.Id))
}

// 通过order_id更新订单信息
func (orderService *OrderService) Update(ctx *gin.Context) interface{} {
	orderService.registerEvent()
	var orderModel models.OrderModel
	byteData, _ := ctx.GetRawData()
	if err := json.Unmarshal(byteData, &orderModel); err != nil {
		return "数据解析异常，请核对：" + err.Error()
	}
	mysqlClient := mysqllib.GetMysqlClient()
	obj := make(map[string]interface{})
	if err := json.Unmarshal(byteData, &obj); err != nil {
		return "数据解析异常，请核对：" + err.Error()
	}
	//多组批量更新
	result := mysqlClient.Model(&models.OrderModel{}).Where("order_id = ?", obj["order_id"]).Updates(obj)

	// 指定字段更新。使用 Struct 进行 Select（会 select 零值的字段）
	//result :=mysqlClient.Model(&orderModel).Select("Name", "Age").Updates(User{Name: "new_name", Age: 0})

	// Select 所有字段（查询包括零值字段的所有字段）
	//db.Model(&user).Select("Name", "Age").Updates(User{Name: "new_name", Age: 0})
	//mysqlClient.Model(&orderModel).Select("*").Update(models.OrderModel{
	//	Id:               0,
	//	OrderId:          "",
	//	PlatformCode:     "",
	//	AccountId:        "",
	//	OrderStatus:      "",
	//	ShipName:         "",
	//	ShipStreet1:      "",
	//	ShipCountry:      "",
	//	ShipCityName:     "",
	//	ShipCode:         "",
	//	ShipPhone:        "",
	//	MiddleCreateTime: utils.LocalTime{},
	//})

	//发送mq通知程序，更新es信息，引用传递
	orderService.eventManager.Trigger(constants.EventOrderChange, obj["order_id"].(string))
	return "更新订单成功，id为：" + strconv.Itoa(int(result.RowsAffected))
}

// 通过order_id删除订单
func (orderService *OrderService) Delete(ctx *gin.Context) interface{} {
	orderService.registerEvent()
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
	//发送mq通知程序，更新es信息，引用传递
	orderService.eventManager.Trigger(constants.EventOrderChange, orderId)
	return "删除订单成功:" + strconv.Itoa(int(result.RowsAffected))
}

// 注册绑定事件
func (orderService *OrderService) registerEvent() {
	orderService.eventManager.Bind(constants.EventOrderChange, &events.OrderEventHandler{})
}
