package graph

type Iter struct {
	g       GraphI
	iter    []int
	reverse bool // iter from 1st layout
	gap     bool // gap for loop
}

func NewIter(g GraphI, reverse bool) *Iter {
	if g == nil || g.Len() == 0 {
		return nil
	}
	for i := 0; i < g.Len(); i++ {
		l := g.Get(i)
		if l == nil || l.Len() == 0 {
			return nil
		}
	}
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
		layer := iter.g.Get(i)
		nodes = append(nodes, layer.Get(it%layer.Len()))
	}
	return nodes
}
