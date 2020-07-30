package suture

import (
	"errors"
)

var (
	ErrComplete = errors.New("complete")
	ErrAbort    = errors.New("abort")
)
