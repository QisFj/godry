package informer

import "context"

// Interface is used to access remote resources.
// may implmented by HTTP API or MySQL, etc.
type Interface[ObjectContent any] interface {
	Create(ctx context.Context, object Object[ObjectContent]) (Object[ObjectContent], error)

	// List return all objects
	List(ctx context.Context) ([]Object[ObjectContent], error)
	Get(ctx context.Context, name string) (Object[ObjectContent], error)

	Update(ctx context.Context, object Object[ObjectContent]) (Object[ObjectContent], error)

	Delete(ctx context.Context, name string) error
}

// todo[maybe]: Watchable Interface
