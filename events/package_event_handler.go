package events

type PackageEventHandler struct {
}

func (packageEventHandler *PackageEventHandler) Handle(params interface{}) {
	//orderId := params.(string)
	//rabbitMQ := services.NewRabbitMQ()
	//orderMap := map[string]string{"order_id": orderId}
	//messageJson, _ := json.Marshal(orderMap)
	//rabbitMQ.SendMessage(string(messageJson), services.ExchangePackageChange, services.QueuePackageChange)
}
