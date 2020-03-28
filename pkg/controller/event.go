package controller

type EventType int

const (
	EventAdded EventType = iota
	EventUpdated
	EventDeleted
)

type Event struct {
	Type EventType
	Object interface{}
}



type EventsHook interface {
	Hook
	GetEventChan() <- chan Event
}

type eventsHook struct {
	events chan Event
}

func NewEventsHook(channelSize int)EventsHook{
	return &eventsHook{events: make(chan Event, channelSize),}
}

func(e *eventsHook)GetEventChan() <- chan Event {
	return e.events
}

func(e *eventsHook)OnAdd(object interface{}){
	e.events <- Event{
		Type:   EventAdded,
		Object: object,
	}
}


func(e *eventsHook)OnUpdate(object interface{}){
	e.events <- Event{
		Type:   EventUpdated,
		Object: object,
	}
}

func (e *eventsHook)OnDelete(object interface{}) {
	e.events <- Event{
		Type:   EventDeleted,
		Object: object,
	}
}
