package run

import "time"

// EachUntil run f each dur, until stopCh return
func EachUntil(f func(), dur time.Duration, stopCh <-chan struct{}) {
	ticker := time.NewTicker(dur)
	defer ticker.Stop()

	for {
		select {
		case <-stopCh:
			return
		case <-ticker.C:
			f()
		}
	}
}

// EachUntilImmediately like EachUntil
// but run f once immediately, unless stopCh already return
func EachUntilImmediately(f func(), dur time.Duration, stopCh <-chan struct{}) {
	select {
	case <-stopCh:
		return
	default:
		f()
	}
	EachUntil(f, dur, stopCh)
}
