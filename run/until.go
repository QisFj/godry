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

// CheckUntil run f each dur, until f return true or stopCh return
// when stop because f return true, CheckUntil return true
// when stop because stopCh return, CheckUntil return false
func CheckUntil(f func() bool, dur time.Duration, stopCh <-chan struct{}) bool {
	ticker := time.NewTicker(dur)
	defer ticker.Stop()

	for {
		select {
		case <-stopCh:
			return false
		case <-ticker.C:
			if f() {
				return true
			}
		}
	}
}

// CheckUntilImmediately is an immediately version of CheckUntil
// CheckUntilImmediatelylike EachUntilImmediately to EachUntil
func CheckUntilImmediately(f func() bool, dur time.Duration, stopCh <-chan struct{}) bool {
	select {
	case <-stopCh:
		return false
	default:
		if f() {
			return true
		}
	}
	return CheckUntil(f, dur, stopCh)
}
