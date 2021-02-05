package iacloaderror

// LoadError is generated when any error related to iac load happens
type LoadError struct {
	ErrMessage string
	Err        error
}

func (e *LoadError) Error() string {
	return e.ErrMessage
}
