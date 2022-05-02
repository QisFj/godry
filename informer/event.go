package informer

import (
	"github.com/QisFj/godry/slice"
)

type EventType string

const (
	CreateEvent EventType = "CREATE"
	UpdateEvent EventType = "UPDATE"
	DeleteEvent EventType = "DELETE"
)

type Event[ObjectContent any] struct {
	Type EventType

	OldObject *Object[ObjectContent] // nil if create
	NewObject *Object[ObjectContent] // nil if delete
}

type EventDispatcher[ObjectContent any] interface {
	Run(stopCh <-chan struct{})
	Push(event Event[ObjectContent])

	// should never call AddHandler and deleteHandler after Run
	AddHandler(handler EventHandler[ObjectContent]) (deleteHandler func())
}

type EventHandler[ObjectContent any] interface {
	Handle(event Event[ObjectContent])
}

type EventHandlerFunc[ObjectContent any] func(event Event[ObjectContent])

func (ehf EventHandlerFunc[ObjectContent]) Handle(event Event[ObjectContent]) {
	if ehf != nil {
		ehf(event)
	}
}

type TypedEventHandler[ObjectContent any] struct {
	OnCreate EventHandlerFunc[ObjectContent]
	OnUpdate EventHandlerFunc[ObjectContent]
	OnDelete EventHandlerFunc[ObjectContent]
}

func (teh TypedEventHandler[ObjectContent]) Handle(event Event[ObjectContent]) {
	switch event.Type {
	case CreateEvent:
		teh.OnCreate.Handle(event)
	case UpdateEvent:
		teh.OnUpdate.Handle(event)
	case DeleteEvent:
		teh.OnDelete.Handle(event)
	}
}

func NewEventDispatcher[ObjectContent any]() EventDispatcher[ObjectContent] {
	return &eventDispatcher[ObjectContent]{
		bc: NewBufferedChannel[Event[ObjectContent]](8),
	}
}

type eventDispatcher[ObjectContent any] struct {
	handlers []EventHandler[ObjectContent]

	bc *bufferedChannel[Event[ObjectContent]]
}

func (ed *eventDispatcher[ObjectContent]) Run(stopCh <-chan struct{}) {
	slice.FilterOn(&ed.handlers, func(_ int, h EventHandler[ObjectContent]) bool {
		return h != nil
	})
	go ed.bc.Run()
	defer close(ed.bc.Source())
	for {
		select {
		case <-stopCh:
			return
		case event := <-ed.bc.Sink():
			for _, handler := range ed.handlers {
				handler.Handle(event)
			}
		}
	}
}

func (ed *eventDispatcher[ObjectContent]) Push(event Event[ObjectContent]) {
	ed.bc.Source() <- event
}

func (ed *eventDispatcher[ObjectContent]) AddHandler(handler EventHandler[ObjectContent]) (deleteHandler func()) {
	i := len(ed.handlers)
	ed.handlers = append(ed.handlers, handler)
	return func() { ed.handlers[i] = nil }
}
