package config

import (
	"time"

	"github.com/QisFj/godry/run"
)

type Reloadable interface {
	Reload() error
}

// ReloadEach call Reload for reloadable each dur, ignore error, util stopCh be closed
func ReloadEach(reloadable Reloadable, dur time.Duration, stopCh <-chan struct{}) {
	run.EachUntil(func() {
		_ = reloadable.Reload() // ignore error
	}, dur, stopCh)
}
