package run

import (
	"sync"

	"github.com/QisFj/godry/channels"
)

// LazyRunner
// call Run when at least one supervisor is added
// stop Run when all supervisors are stopped
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
		lr.runSupervisorChecker()
	}()
	lr.Run(stopCh)
	lr.Locker.Lock()
	defer lr.Locker.Unlock()
	close(lr.stoppedCh)
	lr.stoppedCh = nil
}

func (lr *LazyRunner) runSupervisorChecker() {
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
				break
			}
			stopChs = lr.supervisorStopChs
			lr.Locker.Unlock()
		}

		chosen, _, _ := channels.Read(stopChs)

		lr.Locker.Lock()
		// optimization: swap chosen and last, and truncate
		lr.supervisorStopChs[chosen], lr.supervisorStopChs[len(lr.supervisorStopChs)-1] = lr.supervisorStopChs[len(lr.supervisorStopChs)-1], lr.supervisorStopChs[chosen]
		lr.supervisorStopChs = lr.supervisorStopChs[:len(lr.supervisorStopChs)-1]
		lr.Locker.Unlock()
	}
}

func (lr *LazyRunner) Wait() {
	lr.Locker.RLock()
	stoppedCh := lr.stoppedCh
	lr.Locker.RUnlock()
	if stoppedCh != nil {
		<-stoppedCh
	}
}
