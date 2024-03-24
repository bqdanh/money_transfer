package exceptions

import "fmt"

type InvalidArgumentError struct {
	Field       string
	Description string
	Metadata    map[string]interface{}
}

func NewInvalidArgumentError(field, description string, md map[string]interface{}) *InvalidArgumentError {
	return &InvalidArgumentError{
		Field:       field,
		Description: description,
		Metadata:    md,
	}
}

func (e *InvalidArgumentError) Error() string {
	if e == nil {
		return "nil"
	}
	return fmt.Sprintf("InvalidArgumentError: {%s: %s: %v}", e.Field, e.Description, e.Metadata)
}

func (e *InvalidArgumentError) Is(target error) bool {
	t, ok := target.(*InvalidArgumentError)
	if !ok {
		return false
	}
	return e.Field == t.Field && e.Description == t.Description
}

func (e *InvalidArgumentError) As(target interface{}) bool {
	t, ok := target.(*InvalidArgumentError)
	if !ok {
		return false
	}
	*t = *e
	return true
}
