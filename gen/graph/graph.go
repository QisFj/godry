package graph

// GraphI is a list of LayerIs
type GraphI interface {
	Len() int
	Get(i int) LayerI
}

// LayerI is a list of NodeI
type LayerI interface {
	Len() int
	Get(i int) NodeI
}

type NodeI interface{}

type Graph []LayerI

func (g Graph) Len() int         { return len(g) }
func (g Graph) Get(i int) LayerI { return g[i] }

type Layer []NodeI

func (l Layer) Len() int        { return len(l) }
func (l Layer) Get(i int) NodeI { return l[i] }

type Func struct {
	Length  int // length should be static
	GetFunc func(i int) NodeI
}

func (fn Func) Len() int        { return fn.Length }
func (fn Func) Get(i int) NodeI { return fn.GetFunc(i) }

type NumberRange struct {
	From   int // inclusive
	Length int
}

func (nr NumberRange) Len() int        { return nr.Length }
func (nr NumberRange) Get(i int) NodeI { return nr.From + i }
