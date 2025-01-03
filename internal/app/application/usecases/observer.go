package usecase

import "time"

type EventType string

const InvoiceGenerated EventType = "InvoiceGenerated"

type InvoiceGeneratedEventData struct {
	Amount    float64
	Date      time.Time
	UserEmail string
}

type Event[T any] struct {
	Type EventType
	Date time.Time
	Data T
}

type Observer[T any] struct {
	subscribers map[EventType][]chan<- Event[T]
}

func (o *Observer[T]) Notify(event Event[T]) {
	for _, sub := range o.subscribers[event.Type] {
		sub <- event
	}
}

func (o *Observer[T]) Subscribe(eventType EventType, channel chan<- Event[T]) {
	o.subscribers[eventType] = append(o.subscribers[eventType], channel)
}

func NewObserver[T any]() *Observer[T] {
	return &Observer[T]{subscribers: make(map[EventType][]chan<- Event[T])}
}
