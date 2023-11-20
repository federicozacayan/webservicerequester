package request

// CustomError is a custom error type.
type CustomError struct {
	Message string
}

// Error implements the error interface for CustomError.
func (e *CustomError) Error() string {
	return e.Message
}
