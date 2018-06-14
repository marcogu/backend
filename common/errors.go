package common

type FormValidateError struct {
	Message string
}

func (err FormValidateError) Error() string {
	return err.Message
}

func NewFormValidateError(msg string) *FormValidateError {
	return &FormValidateError{msg}
}
