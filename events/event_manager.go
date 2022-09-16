package events

// 事件管理器
type EventManager struct {
	events map[string]EventHandler
}

// 绑定事件
func (eventManager *EventManager) Bind(eventName string, eventHandle EventHandler) *EventManager {
	isBind := false
	for name, _ := range eventManager.events {
		if eventName == name {
			isBind = true
		}
	}
	if !isBind {
		if eventManager.events == nil {
			eventManager.events = make(map[string]EventHandler)
		}
		eventManager.events[eventName] = eventHandle
	}
	return eventManager
}

// 触发事件
func (eventManager *EventManager) Trigger(eventName string, params interface{}) {
	for name, event := range eventManager.events {
		if eventName == name {
			event.Handle(params)
		}
	}
}
