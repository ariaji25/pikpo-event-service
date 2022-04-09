package events

import (
	"context"
)

type Events struct {
	UnimplementedEventServer
}

var handler = EventHandler{}

func (e *Events) CreateEvent(ctx context.Context, in *EventItem) (*EventItem, error) {
	return handler.CreateEventHandler(in)
}
func (e *Events) GetEvents(ctx context.Context, in *GetEventRequest) (*GetEventResponse, error) {
	return handler.GetEventsHandler(in)
}
func (e *Events) UpdateEvent(ctx context.Context, in *EventItem) (*EventItem, error) {
	return handler.UpdateEventHandler(in)
}
func (e *Events) DeleteEvent(ctx context.Context, in *EventByID) (*DeleteResponse, error) {
	return handler.DeleteEventHandler(in)
}
