package utils

import "fmt"

type LoadBalancerError struct {
	Code    int
	Message string
	Err     error
}

func (e *loadbalancerError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("LoadBalancer Error %d: %s - %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("LoadBalancer Error %d: %s", e.Code, e.Message)
}

func NewError(code int, message string, err error) *loadbalancerError {
	return &LoadBalancerError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Common error codes
const (
	ErrNoBackendsAvailable = 1001
	ErrBackendUnreachable  = 1002
	ErrInvalidConfig       = 1003
)