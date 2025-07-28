package exceptions

type BusinessError struct {
	Code       int
	Message    string
	HttpStatus int
}

func (e *BusinessError) Error() string {
	return e.Message
}
