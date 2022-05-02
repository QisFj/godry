package informer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/QisFj/godry/run"
	"github.com/go-playground/validator/v10"
)

type Informer[ObjectContent any] interface {
	Run(ctx context.Context)

	// LastSyncTime return last successful sync time, update after List returned
	LastSyncTime() time.Time

	Interface() Interface[ObjectContent]

	Lister[ObjectContent]
	// todo[maybe]: Indexer
	AddHandler(handler EventHandler[ObjectContent]) (deleteHandler func())
}

type Config struct {
	ResyncPeriod time.Duration `validate:"required"`

	// OnListError will be called when List error
	// it will block list
	// can be nil, will ignore list error
	OnListError func(err error, continuousCount int)

	// Clock use to mock time-related functions
	// can be nil, will use realClock
	Clock Clock
}

var configValidator = validator.New()

func New[ObjectContent any](config Config, intf Interface[ObjectContent]) (Informer[ObjectContent], error) {
	if err := configValidator.Struct(config); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}
	if config.OnListError == nil {
		config.OnListError = func(err error, continuousCount int) {} // empty function to drop error
	}
	if config.Clock == nil {
		config.Clock = realClock{}
	}
	return &informer[ObjectContent]{
		config:          config,
		intf:            intf,
		eventDispatcher: NewEventDispatcher[ObjectContent](),
	}, nil
}

// informer is implmentation of Informer
type informer[ObjectContent any] struct {
	config Config

	// rw protect:
	// - lastSyncTime
	// - objectes
	rw sync.RWMutex

	lastSyncTime time.Time
	objects      map[string]*Object[ObjectContent]

	listContinuousErrorCount int

	intf            Interface[ObjectContent]
	eventDispatcher EventDispatcher[ObjectContent]

	Informer[ObjectContent]
}

func (inf *informer[ObjectContent]) Run(ctx context.Context) {
	var wg sync.WaitGroup

	stopEventDispatcherCh := make(chan struct{})
	wg.Add(1)
	go func() {
		defer wg.Done()
		run.EachUntilImmediately(func() {
			inf.sync(ctx)
		}, inf.config.ResyncPeriod, ctx.Done())
		close(stopEventDispatcherCh)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		// should stop eventDispatcher after sync loop stopped
		// to avoid write new event after eventDispatcher stopped
		inf.eventDispatcher.Run(stopEventDispatcherCh)
	}()

	wg.Wait()

}

func (inf *informer[ObjectContent]) sync(ctx context.Context) {
	objs, err := inf.intf.List(ctx)
	if err != nil {
		inf.listContinuousErrorCount++
		inf.config.OnListError(err, inf.listContinuousErrorCount)
		return
	}
	newObjects := make(map[string]*Object[ObjectContent])
	for i := range objs {
		newObjects[objs[i].Name] = &objs[i]
	}
	inf.listContinuousErrorCount = 0

	inf.rw.Lock()
	inf.lastSyncTime = inf.config.Clock.Now()
	oldObjects := inf.objects
	inf.objects = newObjects
	inf.rw.Unlock()

	// push events
	for name, newObj := range newObjects {
		oldObj := oldObjects[name]
		if oldObj == nil {
			inf.eventDispatcher.Push(Event[ObjectContent]{
				Type:      CreateEvent,
				NewObject: newObj,
			})
			continue
		}
		if oldObj.ResourceVersion != newObj.ResourceVersion {
			inf.eventDispatcher.Push(Event[ObjectContent]{
				Type:      UpdateEvent,
				OldObject: oldObj,
				NewObject: newObj,
			})
			continue
		}
	}
	for name, oldObj := range oldObjects {
		if newObj := newObjects[name]; newObj == nil {
			inf.eventDispatcher.Push(Event[ObjectContent]{
				Type:      DeleteEvent,
				OldObject: oldObj,
			})
		}
	}
}

func (inf *informer[ObjectContent]) List() []*Object[ObjectContent] {
	inf.rw.RLock()
	defer inf.rw.RUnlock()
	objects := make([]*Object[ObjectContent], 0, len(inf.objects))
	for _, obj := range inf.objects {
		objects = append(objects, obj)
	}
	return objects
}

func (inf *informer[ObjectContent]) Get(name string) *Object[ObjectContent] {
	inf.rw.RLock()
	defer inf.rw.RUnlock()
	return inf.objects[name]
}

func (inf *informer[ObjectContent]) AddHandler(eh EventHandler[ObjectContent]) (removeHandler func()) {
	return inf.eventDispatcher.AddHandler(eh)
}

func (inf *informer[ObjectContent]) LastSyncTime() time.Time {
	inf.rw.RLock()
	defer inf.rw.RUnlock()
	return inf.lastSyncTime
}

const firstSyncedCheckPeriod = 100 * time.Millisecond

func WaitFirstSyncedWithTimeout[ObjectContent any](stopCh <-chan struct{}, timeout time.Duration, informer Informer[ObjectContent]) bool {
	combinedStopCh := make(chan struct{})
	timer := time.NewTimer(timeout)
	go func() {
		select {
		case <-stopCh:
		case <-timer.C:
		}
		close(combinedStopCh)
		timer.Stop() // expiry timer, so it can GC before timeout
	}()
	return WaitFirstSynced(combinedStopCh, informer)
}

func WaitFirstSynced[ObjectContent any](stopCh <-chan struct{}, informer Informer[ObjectContent]) bool {
	return run.CheckUntilImmediately(func() bool {
		return !informer.LastSyncTime().IsZero()
	}, firstSyncedCheckPeriod, stopCh)
}
