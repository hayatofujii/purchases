package utils

import "fmt"

type HTTPError struct {
	StatusCode int
	Err        error
}

func (r HTTPError) Error() string {
	return fmt.Sprintf("%v", r.Err)
}
