package informer

import (
	"context"
	"sync"
	"time"

	"github.com/QisFj/godry/run"
)

type Informer[ObjectContent any] interface {
	Run(ctx context.Context) error

	HasSynced() bool

	Lister[ObjectContent]

	// todo[maybe]: Indexer
	AddEventandler(handler EventHandler[ObjectContent]) (deleteHandler func())
}

func New[ObjectContent any](lw ListAndWatch[ObjectContent]) (Informer[ObjectContent], error) {
	return &informer[ObjectContent]{
		lw:               lw,
		unhandledChanges: make(map[string]changeLog[ObjectContent]),
		eventDispatcher:  NewEventDispatcher[ObjectContent](),
	}, nil
}

// informer is implmentation of Informer
type informer[ObjectContent any] struct {
	lw ListAndWatch[ObjectContent]

	// rw protect:
	// - hasSynced
	// - unhandledChanges
	// - objectes
	rw sync.RWMutex

	hasSynced        bool
	unhandledChanges map[string]changeLog[ObjectContent]
	objects          map[string]*Object[ObjectContent]

	eventDispatcher EventDispatcher[ObjectContent]
}

type changeLog[ObjectContent any] struct {
	changeType string // set, del
	meta       ObjectMeta
	content    ObjectContent // only make sense when changeType == "set"
}

func (inf *informer[ObjectContent]) Run(ctx context.Context) error {
	errCh := make(chan error, 1)
	occurErr := func(err error) {
		select {
		case errCh <- err:
		default:
		}
	}

	var wg sync.WaitGroup

	watchingCh := make(chan struct{})
	watchStopped := make(chan struct{})
	wg.Add(1)
	go func() {
		defer wg.Done()
		occurErr(inf.lw.Watch(ctx, Callback[ObjectContent]{
			OnCreate: inf.onSet,
			OnUpdate: inf.onSet,
			OnDelete: inf.onDelete,
		}, watchingCh))
		close(watchStopped)
	}()

	go func() {
		select {
		case <-ctx.Done():
			return
		case <-watchingCh:
		}
		run.CheckUntilImmediately(func() bool {
			objects, err := inf.lw.List(ctx)
			if err != nil {
				// todo[maybe]: log error
				return false
			}
			inf.rw.Lock()
			defer inf.rw.Unlock()
			for _, obj := range objects {
				inf._recordSetChange(obj)
			}
			inf._handleChanges()
			inf.hasSynced = true
			return true
		}, 3*time.Second, watchStopped) // stop retry if watchStopped
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		// should stop eventDispatcher after sync loop stopped
		// to avoid write new event after eventDispatcher stopped
		inf.eventDispatcher.Run(watchStopped)
	}()

	wg.Wait()

	select {
	case err := <-errCh:
		return err
	default:
		return nil
	}
}

func (inf *informer[ObjectContent]) onSet(obj Object[ObjectContent]) {
	inf.rw.Lock()
	defer inf.rw.Unlock()

	if !inf.hasSynced {
		inf._recordSetChange(obj)
		return
	}

	oldObj := inf.objects[obj.Name]
	if oldObj == nil {
		// create
		inf.objects[obj.Name] = &obj
		inf.eventDispatcher.Push(Event[ObjectContent]{
			Type:      CreateEvent,
			NewObject: &obj,
		})
		return
	}
	if oldObj.ResourceVersion < obj.ResourceVersion {
		// update
		inf.objects[obj.Name] = &obj
		inf.eventDispatcher.Push(Event[ObjectContent]{
			Type:      UpdateEvent,
			OldObject: oldObj,
			NewObject: &obj,
		})
	}
	// else, drop
}

func (inf *informer[ObjectContent]) onDelete(meta ObjectMeta) {
	inf.rw.Lock()
	defer inf.rw.Unlock()
	if !inf.hasSynced {
		inf._recordDelChange(meta)
		return
	}

	oldObj := inf.objects[meta.Name]
	if oldObj == nil {
		// drop
		return
	}
	// delete
	delete(inf.objects, meta.Name)
	inf.eventDispatcher.Push(Event[ObjectContent]{
		Type:      DeleteEvent,
		OldObject: oldObj,
	})
}

// _recordSetChange should be called under rw.Lock()
func (inf *informer[ObjectContent]) _recordSetChange(obj Object[ObjectContent]) {
	change, ok := inf.unhandledChanges[obj.Name]
	if ok && change.meta.ResourceVersion >= obj.ResourceVersion {
		// keep the origin change
		return
	}
	inf.unhandledChanges[obj.Name] = changeLog[ObjectContent]{
		changeType: "set",
		meta:       obj.ObjectMeta,
		content:    obj.Content,
	}
}

// _recordDelChange should be called under rw.Lock()
func (inf *informer[ObjectContent]) _recordDelChange(meta ObjectMeta) {
	change, ok := inf.unhandledChanges[meta.Name]
	if ok && change.meta.ResourceVersion > meta.ResourceVersion {
		// keep the origin change
		return
	}
	inf.unhandledChanges[meta.Name] = changeLog[ObjectContent]{
		changeType: "del",
		meta:       meta,
	}
}

// _handleChanges should be called under rw.Lock()
func (inf *informer[ObjectContent]) _handleChanges() {
	inf.objects = map[string]*Object[ObjectContent]{}
	for name, change := range inf.unhandledChanges {
		if change.changeType == "set" {
			obj := &Object[ObjectContent]{
				ObjectMeta: change.meta,
				Content:    change.content,
			}
			inf.objects[name] = obj
			inf.eventDispatcher.Push(Event[ObjectContent]{
				Type:      CreateEvent,
				NewObject: obj,
			})
		}
		// ignore delete
	}
}

func (inf *informer[ObjectContent]) Get(name string) *Object[ObjectContent] {
	inf.rw.RLock()
	defer inf.rw.RUnlock()
	return inf.objects[name]
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

func (inf *informer[ObjectContent]) AddEventandler(eh EventHandler[ObjectContent]) (removeHandler func()) {
	return inf.eventDispatcher.AddHandler(eh)
}

func (inf *informer[ObjectContent]) HasSynced() bool {
	inf.rw.RLock()
	defer inf.rw.RUnlock()
	return inf.hasSynced
}

const checkSyncedPeriod = 100 * time.Millisecond

func WaitSyncedWithTimeout[ObjectContent any](stopCh <-chan struct{}, timeout time.Duration, informer Informer[ObjectContent]) bool {
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
	return WaitSynced(combinedStopCh, informer)
}

func WaitSynced[ObjectContent any](stopCh <-chan struct{}, informer Informer[ObjectContent]) bool {
	return run.CheckUntilImmediately(informer.HasSynced, checkSyncedPeriod, stopCh)
}
