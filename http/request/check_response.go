package request

type CheckResponseBeforeUnmarshalFunc func(statusCode int, body []byte) error

type CheckResponseBeforeUnmarshalFuncs []CheckResponseBeforeUnmarshalFunc

func (fs CheckResponseBeforeUnmarshalFuncs) Check(statusCode int, body []byte) error {
	for _, f := range fs {
		if err := f(statusCode, body); err != nil {
			return err
		}
	}
	return nil
}

type CheckResponseAfterUnmarshalFunc func(statusCode int, v interface{}) error

type CheckResponseAfterUnmarshalFuncs []CheckResponseAfterUnmarshalFunc

func (fs CheckResponseAfterUnmarshalFuncs) Check(statusCode int, v interface{}) error {
	for _, f := range fs {
		if err := f(statusCode, v); err != nil {
			return err
		}
	}
	return nil
}
