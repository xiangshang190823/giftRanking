package errors

type error interface {
	Error() string
}

func New(code string,text string) error {
	return &errorString{code,text}
}
// errorString is a trivial implementation of errors.
type errorString struct {
	code string
	s string
}
func (e *errorString) Error() string {
	return e.s
}