package sort

type Heap struct {
	// Options
	Less func(o1, o2 interface{}) bool // required
	Size int                           // optional, 0 to not limit

	heap []interface{}
}

func (h *Heap) Append(v interface{}) {
	if h.Size != 0 && h.Size == len(h.heap) {
		if h.Less(h.heap[0], v) {
			// no need append
			return
		}
		// replace top and heapify
		h.heap[0] = v
		h.dive(0)
	} else {
		h.heap = append(h.heap, v)
		h.rise(len(h.heap) - 1)
	}
}

func (h *Heap) Dump() []interface{} {
	result := make([]interface{}, len(h.heap))
	for n := len(h.heap) - 1; n >= 0; n-- {
		result[n] = h.heap[0]
		h.heap[0] = h.heap[n]
		h.heap = h.heap[:n]
		h.dive(0)
	}
	return result
}

func (h *Heap) rise(idx int) {
	for idx != 0 {
		if idx == 0 {
			// already risen to the top
			break
		}
		top := (idx - 1) / 2
		if h.Less(h.heap[idx], h.heap[top]) {
			// no need to continue rise
			break
		}
		h.heap[idx], h.heap[top] = h.heap[top], h.heap[idx]
		idx = top
	}
}

func (h *Heap) dive(idx int) {
	for {
		l, r := idx*2+1, idx*2+2
		if l >= len(h.heap) {
			// already dived to the bottom
			break
		}
		greater := l
		if r < len(h.heap) && h.Less(h.heap[l], h.heap[r]) {
			// right exist, and greater
			greater = r
		}
		if h.Less(h.heap[greater], h.heap[idx]) {
			// no need to continue dive
			break
		}
		h.heap[idx], h.heap[greater] = h.heap[greater], h.heap[idx]
		idx = greater
	}
}
