package informer

type ObjectMeta struct {
	Name            string
	ResourceVersion string
}

type Object[Content any] struct {
	ObjectMeta

	Content Content
}
