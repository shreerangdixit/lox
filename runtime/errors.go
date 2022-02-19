package runtime

// Loop control flow exit due to `break;`
type BreakError struct{}

func NewBreakError() error {
	return BreakError{}
}

func (e BreakError) Error() string { return "break" }

// Function control flow exit due to `return <exp>;`
type ReturnError struct {
	Value Object
}

func NewReturnError(value Object) error {
	return ReturnError{
		Value: value,
	}
}

func (e ReturnError) Error() string { return "return" }
