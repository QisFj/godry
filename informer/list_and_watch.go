package informer

import (
	"container/list"
	"context"
	"sync"
	"time"

	"github.com/QisFj/godry/run"
)

type ListAndWatch[ObjectContent any] interface {
	List(ctx context.Context) ([]Object[ObjectContent], error)
	Watch(ctx context.Context, cb Callback[ObjectContent], watchingNotifyCh chan<- struct{}) error
}

type Callback[ObjectContent any] struct {
	OnCreate func(obj Object[ObjectContent])
	OnUpdate func(obj Object[ObjectContent])
	OnDelete func(meta ObjectMeta)
}

// todo[must]: Create, Update, Delete

type PollListAndWatch[ObjectContent any] struct {
	ListFunc func(ctx context.Context) ([]Object[ObjectContent], error) // required

	PollPeriod time.Duration // required

	// OnListError will be called when List error
	// it will block list
	// can be nil, will ignore list error
	OnListError func(err error, continuousCount int)

	rw         sync.RWMutex // protect: objects, lazyRunner, cblist, cb
	objects    map[string]*Object[ObjectContent]
	lazyRunner *run.LazyRunner

	cblist *list.List
	cb     Callback[ObjectContent] // build from cblist

	listContinuousErrorCount int
}

func (p *PollListAndWatch[ObjectContent]) List(ctx context.Context) ([]Object[ObjectContent], error) {
	return p.ListFunc(ctx)
}

func (p *PollListAndWatch[ObjectContent]) Watch(ctx context.Context, cb Callback[ObjectContent], watchingNotifyCh chan<- struct{}) error {
	p.initLazyRunner()
	// add listener
	p.rw.Lock()
	if p.cblist == nil {
		p.cblist = list.New()
	}
	elem := p.cblist.PushBack(cb)
	p._buildCallback()
	p.rw.Unlock()

	p.lazyRunner.AddSupervisor(ctx.Done())

	if watchingNotifyCh != nil {
		close(watchingNotifyCh)
	}

	<-ctx.Done()

	p.rw.Lock()
	p.cblist.Remove(elem)
	p._buildCallback()
	p.rw.Unlock()
	return nil
}

// _buildCallback should be called under rw.Lock
func (p *PollListAndWatch[ObjectContent]) _buildCallback() {
	cbs := make([]Callback[ObjectContent], 0, p.cblist.Len())
	for e := p.cblist.Front(); e != nil; e = e.Next() {
		cbs = append(cbs, e.Value.(Callback[ObjectContent]))
	}
	p.cb = Callback[ObjectContent]{
		OnCreate: func(obj Object[ObjectContent]) {
			for _, cb := range cbs {
				cb.OnCreate(obj)
			}
		},
		OnUpdate: func(obj Object[ObjectContent]) {
			for _, cb := range cbs {
				cb.OnUpdate(obj)
			}
		},
		OnDelete: func(meta ObjectMeta) {
			for _, cb := range cbs {
				cb.OnDelete(meta)
			}
		},
	}
}

func (p *PollListAndWatch[ObjectContent]) initLazyRunner() {
	p.rw.Lock()
	defer p.rw.Unlock()
	if p.lazyRunner != nil {
		return
	}
	p.lazyRunner = &run.LazyRunner{
		Run: func(stopCh <-chan struct{}) {
			ctx, cancel := context.WithCancel(context.Background())
			go func() {
				<-stopCh
				cancel()
			}()
			run.EachUntilImmediately(func() {
				p.sync(ctx)
			}, p.PollPeriod, ctx.Done())
		},
		Locker: &p.rw,
	}
}

func (p *PollListAndWatch[ObjectContent]) sync(ctx context.Context) {
	objs, err := p.ListFunc(ctx)
	if err != nil {
		p.listContinuousErrorCount++
		if p.OnListError != nil {
			p.OnListError(err, p.listContinuousErrorCount)
		}
		return
	}
	newObjects := make(map[string]*Object[ObjectContent])
	for i := range objs {
		newObjects[objs[i].Name] = &objs[i]
	}
	p.listContinuousErrorCount = 0

	p.rw.Lock()
	defer p.rw.Unlock()
	oldObjects := p.objects
	p.objects = newObjects

	// push events
	for name, newObj := range newObjects {
		oldObj := oldObjects[name]
		if oldObj == nil {
			p.cb.OnCreate(*newObj)
			continue
		}
		if oldObj.ResourceVersion != newObj.ResourceVersion {
			// todo[maybe]: no need check if equal?
			p.cb.OnUpdate(*newObj)
			continue
		}
	}
	for name, oldObj := range oldObjects {
		if newObj := newObjects[name]; newObj == nil {
			p.cb.OnDelete(oldObj.ObjectMeta)
		}
	}
}

func (p *PollListAndWatch[ObjectContent]) Create(obj Object[ObjectContent]) {
	p.rw.Lock()
	defer p.rw.Unlock()
	oldObj := p.objects[obj.Name]
	if oldObj != nil {
		return
	}
	p.objects[obj.Name] = &obj
	p.cb.OnCreate(obj)
}

func (p *PollListAndWatch[ObjectContent]) Update(obj Object[ObjectContent]) {
	p.rw.Lock()
	defer p.rw.Unlock()
	oldObj := p.objects[obj.Name]
	if oldObj != nil && oldObj.ResourceVersion >= obj.ResourceVersion {
		return
	}
	p.objects[obj.Name] = &obj
	p.cb.OnUpdate(obj)
}

func (p *PollListAndWatch[ObjectContent]) Delete(meta ObjectMeta) {
	p.rw.Lock()
	defer p.rw.Unlock()
	oldObj := p.objects[meta.Name]
	if oldObj == nil {
		return
	}
	delete(p.objects, meta.Name)
	p.cb.OnDelete(meta)
}
