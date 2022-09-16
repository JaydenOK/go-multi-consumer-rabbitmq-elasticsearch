package events

import (
	"app/constants"
	"app/libs/rabbitmqlib"
	"encoding/json"
	"fmt"
)

type OrderEventHandler struct {
}

func (orderEventHandle *OrderEventHandler) Handle(params interface{}) {
	fmt.Println("发送消息：", params)
	orderId := params.(string)
	rabbitMQ := rabbitmqlib.NewRabbitMQ()
	orderMap := map[string]string{"order_id": orderId}
	messageJson, _ := json.Marshal(orderMap)
	if err := rabbitMQ.SendMessage(string(messageJson), constants.EventOrderChange, constants.QueueOrderChange); err != nil {
		fmt.Println("发送消息失败：", err.Error())
	}
}
