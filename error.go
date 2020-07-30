package suture

import (
	"errors"
)

var (
	ErrComplete = errors.New("complete")
	ErrTearDown = errors.New("tear down")
)
