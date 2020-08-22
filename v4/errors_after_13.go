package suture

import "errors"

// +build go1.13
func isErr(err error, target error) bool {
	return errors.Is(err, target)
}
