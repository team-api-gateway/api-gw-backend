package domain

import "fmt"

type HttpError struct {
	Err    error
	Status int
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("%v", e.Err)
}
