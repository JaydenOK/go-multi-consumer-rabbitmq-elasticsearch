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
