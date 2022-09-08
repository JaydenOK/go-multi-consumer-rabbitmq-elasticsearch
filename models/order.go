package models

type Order struct {
	Id                  uint   `json:"id"`
	OrderId             string `json:"order_id"`
	PlatformCode        string `json:"platform_code"`
	AccountId           string `json:"account_id"`
	OrderStatus         string `json:"order_status"`
	ShipName            string `json:"ship_name"`
	ShipStreet1         string `json:"ship_street1"`
	ShipCountry         string `json:"ship_country"`
	ShipCityName        string `json:"ship_city_name"`
	ShipStateOrProvince string `json:"ship_stateorprovince"`
	MiddleCreateTime    string `json:"middle_create_time"`
}

func (order *Order) TableName() string {
	return "yb_order"
}


