package models

import "app/utils"

type OrderModel struct {
	Id               uint            `json:"id"`
	OrderId          string          `json:"order_id"`
	PlatformCode     string          `json:"platform_code"`
	AccountId        string          `json:"account_id"`
	OrderStatus      string          `json:"order_status"`
	ShipName         string          `json:"ship_name"`
	ShipStreet1      string          `json:"ship_street1"`
	ShipCountry      string          `json:"ship_country"`
	ShipCityName     string          `json:"ship_city_name"`
	ShipCode         string          `json:"ship_code"`
	ShipPhone        string          `json:"ship_phone"`
	TotalPrice       float64         `json:"total_price"`
	MiddleCreateTime utils.LocalTime `json:"middle_create_time"` //utils.LocalTime： 实现MarshalJSON接口，格式化数据
}

func (orderModel *OrderModel) TableName() string {
	return "yb_order"
}

//es设置mapping结构：
//{"mappings":{"properties":{"account_id":{"fields":{"keyword":{"ignore_above":256,"type":"keyword"}},"type":"keyword"},"id":{"type":"long"},"middle_create_time":{"fields":{"keyword":{"ignore_above":256,"type":"keyword"}},"type":"date","format":"yyyy-MM-dd HH:mm:ss"},"order_id":{"fields":{"keyword":{"ignore_above":256,"type":"keyword"}},"type":"text"},"order_status":{"fields":{"keyword":{"ignore_above":256,"type":"keyword"}},"type":"keyword"},"platform_code":{"fields":{"keyword":{"ignore_above":256,"type":"keyword"}},"type":"keyword"},"ship_city_name":{"fields":{"keyword":{"ignore_above":256,"type":"keyword"}},"type":"keyword"},"ship_code":{"fields":{"keyword":{"ignore_above":256,"type":"keyword"}},"type":"keyword"},"ship_country":{"fields":{"keyword":{"ignore_above":256,"type":"keyword"}},"type":"keyword"},"ship_name":{"fields":{"keyword":{"ignore_above":256,"type":"keyword"}},"type":"text"},"ship_phone":{"fields":{"keyword":{"ignore_above":256,"type":"keyword"}},"type":"keyword"},"ship_street1":{"fields":{"keyword":{"ignore_above":256,"type":"keyword"}},"type":"text"},"total_price":{"type":"float"}}}}
