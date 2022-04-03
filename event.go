package fsm

type EventCode string

type Event interface {
	GetEventData() interface{}
	Code() EventCode
}

type defaultEvent struct {
	code      EventCode
	eventData interface{}
}

// NewEvent event with optional data
func NewEvent(eventCode EventCode, eventData ...interface{}) Event {
	if eventCode == "" {
		return nil
	}
	event := &defaultEvent{
		code: eventCode,
	}
	if len(eventData) > 0 {
		event.eventData = eventData[0]
	}
	return event
}

func (d *defaultEvent) GetEventData() interface{} {
	return d.eventData
}

func (d *defaultEvent) Code() EventCode {
	return d.code
}
