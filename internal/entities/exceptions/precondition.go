package exceptions

import "fmt"

type PreconditionType string

const (
	PreconditionTypeUserDuplicatedUserName = PreconditionType("user-name-duplicated")
)

type PreconditionError struct {
	// The type of PreconditionFailure, example: AccountPrecondition
	Type        PreconditionType
	Subject     string
	Description string
	Metadata    map[string]interface{}
}

func NewPreconditionError(ptype PreconditionType, subject string, description string, md map[string]interface{}) *PreconditionError {
	return &PreconditionError{
		Type:        ptype,
		Subject:     subject,
		Description: description,
		Metadata:    md,
	}
}

func (e *PreconditionError) Error() string {
	if e == nil {
		return "nil"
	}
	return fmt.Sprintf("PreconditionError: {type: %s, subject: %s, description: %s}", e.Type, e.Subject, e.Description)
}

func (e *PreconditionError) Is(target error) bool {
	t, ok := target.(*PreconditionError)
	if !ok {
		return false
	}
	return e.Type == t.Type && e.Subject == t.Subject
}

func (e *PreconditionError) As(target interface{}) bool {
	t, ok := target.(*PreconditionError)
	if !ok {
		return false
	}
	*t = *e
	return true
}
