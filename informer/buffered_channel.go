package informer

// bufferedChannel is like a infinite channel, after run
// it stopless read from sourceCh, and write to sinkCh, if write blocked, store into a list
// all stored value will be write to sinkCh once it not blocked
// when sourceCh is closed
// - bufferedChannel will stop run
// - and buffered value will droped
// - and sinkCh will be closed too
type bufferedChannel[V any] struct {
	sourceCh chan V
	sinkCh   chan V

	buffer *infiniteRingBuffer[V]
}

func NewBufferedChannel[V any](initBufferSize int) *bufferedChannel[V] {
	return &bufferedChannel[V]{
		sourceCh: make(chan V),
		sinkCh:   make(chan V),
		buffer:   newInfiniteRingBuffer[V](initBufferSize),
	}
}

func (cb *bufferedChannel[V]) Run() {
	// inspired by: https://github.com/kubernetes/client-go/blob/f6ce18ae578c8cca64d14ab9687824d9e1305a67/tools/cache/shared_informer.go#L736
	defer close(cb.sinkCh)

	var (
		sinkCh chan<- V
		v      V
	)

	for {
		select {
		case newV, ok := <-cb.sourceCh:
			if !ok {
				// sourceCh is closed
				return
			}
			if sinkCh == nil {
				v = newV
				sinkCh = cb.sinkCh
			} else {
				cb.buffer.append(newV)
			}
		case sinkCh <- v:
			var ok bool
			v, ok = cb.buffer.pop() // load next need write value
			if !ok {
				// no buffered value need to write
				sinkCh = nil
			}
		}

	}
}

func (cb *bufferedChannel[V]) Source() chan<- V {
	return cb.sourceCh
}

func (cb *bufferedChannel[V]) Sink() chan V {
	return cb.sinkCh
}

type infiniteRingBuffer[V any] struct {
	ring     []V
	startIdx int
	length   int

	next *infiniteRingBuffer[V]
}

func newInfiniteRingBuffer[V any](initSize int) *infiniteRingBuffer[V] {
	if initSize <= 0 {
		initSize = 64
	}
	return &infiniteRingBuffer[V]{
		ring: make([]V, initSize),
	}
}

func (irb *infiniteRingBuffer[V]) pop() (v V, ok bool) {
	if irb.length == 0 {
		if irb.next != nil {
			*irb = *irb.next
			return irb.pop()
		}
		ok = false
		return
	}
	idx := irb.startIdx

	irb.startIdx = (irb.startIdx + 1) % len(irb.ring)
	irb.length--

	return irb.ring[idx], true
}

func (irb *infiniteRingBuffer[V]) append(v V) {
	if irb.next != nil {
		irb.next.append(v)
		return
	}
	if irb.length == len(irb.ring) {
		irb.next = newInfiniteRingBuffer[V](len(irb.ring) * 2)
		irb.next.append(v)
		return
	}
	idx := (irb.startIdx + irb.length) % len(irb.ring)
	irb.ring[idx] = v
	irb.length++
}
