// +build go1.13

package suture

import "errors"

func isErr(err error, target error) bool {
	return errors.Is(err, target)
}
