package suture

// +build !go1.13

func isErr(err error, target error) bool {
	return err == target
}
