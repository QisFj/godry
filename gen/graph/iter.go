package graph

type Iter struct {
	g       Interface
	iter    []int
	reverse bool // iter from 1st layout
	gap     bool // gap for loop
}

func NewIter(g Interface, reverse bool) *Iter {
	return &Iter{g: g, iter: make([]int, g.Len()), gap: true, reverse: reverse}
}

func (iter *Iter) Next() bool {
	if iter == nil {
		return false
	}
	if iter.gap {
		iter.gap = false
		return true
	}
	for i := range iter.iter {
		if !iter.reverse {
			i = len(iter.iter) - 1 - i
		}
		iter.iter[i]++
		if iter.iter[i]%iter.g.Get(i).Len() != 0 {
			return true
		}
	}
	iter.gap = true
	return false
}

func (iter *Iter) Get() []NodeI {
	nodes := make([]NodeI, 0, len(iter.iter))
	for i, it := range iter.iter {
		nodes = append(nodes, iter.g.Get(i).Get(it%iter.g.Get(i).Len()))
	}
	return nodes
}
