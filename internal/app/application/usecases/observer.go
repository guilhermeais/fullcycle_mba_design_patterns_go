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

func (m *Observer[T]) Notify(event Event[T]) {
	for _, sub := range m.subscribers[event.Type] {
		sub <- event
	}
}

func (m *Observer[T]) Subscribe(eventType EventType, channel chan<- Event[T]) {
	m.subscribers[eventType] = append(m.subscribers[eventType], channel)
}

func NewObserver[T any]() *Observer[T] {
	return &Observer[T]{subscribers: make(map[EventType][]chan<- Event[T])}
}
