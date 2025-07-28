package exceptions

type SystemError struct {
	Code       int
	Message    string
	HttpStatus int
}

func (e *SystemError) Error() string {
	return e.Message
}
