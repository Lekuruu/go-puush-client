package puush

import "strings"

var (
	PuushErrorInvalidCredentials  PuushError = NewPuushError("Authentication failure", -1, false)
	PuushErrorRequestFailure      PuushError = NewPuushError("Connection error", -2, true)
	PuushErrorChecksumFailure     PuushError = NewPuushError("Checksum error", -3, true)
	PuushErrorInsufficientStorage PuushError = NewPuushError("Insufficient storage", -4, false)

	/* Custom internal errors */
	PuushErrorNotFound PuushError = NewPuushError("Not found", -998, false)
	PuushErrorUnknown  PuushError = NewPuushError("Unknown error", -999, false)
)

type PuushError interface {
	error
	Value() int
	String() string
	ShouldRetry() bool
}

type puushError struct {
	name        string
	value       int
	shouldRetry bool
}

func (e puushError) Error() string {
	return "puush error: " + strings.ToLower(e.name)
}

func (e puushError) Value() int {
	return e.value
}

func (e puushError) String() string {
	return e.name
}

func (e puushError) ShouldRetry() bool {
	return e.shouldRetry
}

func NewPuushError(name string, value int, shouldRetry bool) PuushError {
	return puushError{
		name:        name,
		value:       value,
		shouldRetry: shouldRetry,
	}
}
