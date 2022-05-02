package run

import (
	"reflect"
	"sync"
)

// LazyRunner run Run only when there is at least one supervisor
type LazyRunner struct {
	Run func(stopCh <-chan struct{})

	Locker interface {
		sync.Locker

		RLock()
		RUnlock()
	} // protect supervisorStopChs, running, stoppedCh

	supervisorStopChs []<-chan struct{}
	running           bool
	stoppedCh         chan struct{}
}

func (lr *LazyRunner) AddSupervisor(stopCh <-chan struct{}) {
	lr.Locker.Lock()
	defer lr.Locker.Unlock()
	lr.supervisorStopChs = append(lr.supervisorStopChs, stopCh)
	if !lr.running {
		lr.running = true
		lr.stoppedCh = make(chan struct{})
		go lr.run()
	}
}

func (lr *LazyRunner) run() {
	stopCh := make(chan struct{})
	go func() {
		defer close(stopCh)
		for {
			lr.Locker.RLock()
			stopChs := lr.supervisorStopChs
			lr.Locker.RUnlock()

			if len(stopChs) == 0 {
				// optimization: use rlock to check, and use lock to check and set
				lr.Locker.Lock()
				if len(lr.supervisorStopChs) == 0 {
					lr.running = false
					lr.Locker.Unlock()
					return
				}
				stopChs = lr.supervisorStopChs
				lr.Locker.Unlock()
			}

			chosen, _, _ := selectChannels(stopChs)

			lr.Locker.Lock()
			lr.supervisorStopChs = append(lr.supervisorStopChs[:chosen], lr.supervisorStopChs[chosen+1:]...)
			lr.Locker.Unlock()
		}
	}()
	lr.Run(stopCh)
	lr.Locker.Lock()
	defer lr.Locker.Unlock()
	close(lr.stoppedCh)
	lr.stoppedCh = nil
}

func (lr *LazyRunner) Wait() {
	lr.Locker.RLock()
	stoppedCh := lr.stoppedCh
	lr.Locker.RUnlock()
	if stoppedCh != nil {
		<-stoppedCh
	}
}

func selectChannels[T any](chans []<-chan T) (chosenIndex int, value T, ok bool) {
	cases := make([]reflect.SelectCase, len(chans))
	for i, ch := range chans {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
	}
	var rValue reflect.Value
	chosenIndex, rValue, ok = reflect.Select(cases)
	if ok {
		// not closed, should read value
		value = rValue.Interface().(T)
	}
	return
}
