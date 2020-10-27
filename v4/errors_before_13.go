// +build !go1.13

package suture

func isErr(err error, target error) bool {
	return err == target
}
