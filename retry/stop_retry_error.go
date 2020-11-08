package retry

import "fmt"

func StopRetryWithError(err error) error {
	if err == nil {
		return nil
	}
	return stopRetryError{originError: err}
}

type stopRetryError struct {
	originError error
}

func (err stopRetryError) Error() string {
	return fmt.Sprintf("stop retry with error: %s", err.originError)
}
