package multierr

type ErrWrapper func(err error) error
