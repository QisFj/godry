package retry

import (
	"sync"
)

type results struct {
	// options:
	size int

	// inner used field
	cleanCount int
	results    []Result
	rwMu       sync.RWMutex
}

func (rs *results) changeSize(newSize int) {
	rs.rwMu.Lock()
	defer rs.rwMu.Unlock()
	rs.size = newSize
}

func (rs *results) append(result Result) {
	rs.rwMu.Lock()
	defer rs.rwMu.Unlock()
	rs.results = append(rs.results, result)
	if rs.size > 0 && len(rs.results) > rs.size {
		rs.results = rs.results[1:] // remove first
		rs.cleanCount++
	}
}

func (rs *results) getLatest() Result {
	rs.rwMu.RLock()
	defer rs.rwMu.RUnlock()
	return rs.get(rs.cleanCount + len(rs.results) - 1)
}

func (rs *results) get(index int) Result {
	rs.rwMu.RLock()
	defer rs.rwMu.RUnlock()
	if index < rs.cleanCount || index >= rs.cleanCount+len(rs.results) {
		return Result{}
	}
	return rs.results[index-rs.cleanCount]
}

type Result struct {
	Valid bool
	Error error
}
