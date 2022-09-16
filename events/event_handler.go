package events

type EventHandler interface {
	Handle(params interface{})
}
