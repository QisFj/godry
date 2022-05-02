package informer

type Lister[ObjectContent any] interface {
	// should treat every returned object as read only

	List() []*Object[ObjectContent]
	Get(name string) *Object[ObjectContent] // return nil, if not exist
}
