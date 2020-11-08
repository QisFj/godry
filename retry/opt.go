package retry

import "time"

type Option struct {
	MaxRunTime    uint64 // 0: always return; 1: execute only once, without retry
	RetryInterval time.Duration
	F             Func
	ResultSize    int // non-positive number means unlimited
	ErrorWrapper  ErrorWrapper
}
