package config

import (
	"sync"
	"time"

	"github.com/QisFj/godry/multierr"
)

func (ma *MultiAccessor) ReloadEach(dur time.Duration, stopCh <-chan struct{}) {
	ticker := time.NewTicker(dur)
	for {
		select {
		case <-stopCh:
			ticker.Stop()
			return
		case <-ticker.C:
			_ = ma.Reload() // ignore error
		}
	}
}

func (ma *MultiAccessor) Reload() error {
	wg := sync.WaitGroup{}
	errs := multierr.New(nil, nil)
	ma.rw.RLock()
	for key, a := range ma.aMap {
		wg.Add(1)
		go func(key string, a *Accessor) {
			defer wg.Done()
			if err := a.Reload(); err != nil {
				errs.Appendf("config accessor(key=%s) reload error: %w", key, err)
			}
		}(key, a)
	}
	ma.rw.RUnlock()
	wg.Wait()
	return errs.Error()
}
