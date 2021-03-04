package graph

type Interface interface {
	Len() int
	Get(i int) LayoutI
}

type LayoutI interface {
	Len() int
	Get(i int) NodeI
}

type NodeI interface{}

type Graph []LayoutI

func (g Graph) Len() int { return len(g) }

func (g Graph) Get(i int) LayoutI { return g[i] }

type Layout []NodeI

func (l Layout) Len() int { return len(l) }

func (l Layout) Get(i int) NodeI { return l[i] }
