// +build go1.12

package obytes

// Error is the type helping defining errors as constants.
type Error string

func (e Error) Error() string { return string(e) }

const (
	ErrInvalidInput = Error("invalid input")
	ErrUnexpected = Error("unexpected error")
)
