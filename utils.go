package suture

import (
	"fmt"
	"time"
)

// Convert an interface{}, that may be an error, to a string.
func maybeErrToString(maybeErr interface{}) string {
	if err, ok := maybeErr.(error); ok {
		return err.Error()
	}
	return fmt.Sprintf("%#v", maybeErr)
}

// Return the value, unless it is 0. If that is the case, return defaultvalue.
func floatOrDefault(value, defaultValue float64) float64 {
	if value != 0 {
		return value
	}
	return defaultValue
}

// Return the value, unless it is 0. If that is the case, return defaultvalue.
func durationOrDefault(value, defaultValue time.Duration) time.Duration {
	if value != 0 {
		return value
	}
	return defaultValue
}
