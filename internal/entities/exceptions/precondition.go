package exceptions

import "fmt"

type PreconditionReason string

const (
	SubjectUser           = "user"
	SubjectAccount        = "account"
	SubjectAuthentication = "authentication"

	PreconditionReasonUserDuplicatedUserName = PreconditionReason("user-name-duplicated")
	PreconditionReasonCannotChangeUserID     = PreconditionReason("cannot-change-user-id")
	PreconditionReasonUserNotFound           = PreconditionReason("user-not-found")
	PreconditionReasonPasswordNotMatch       = PreconditionReason("password-not-match")
	PreconditionReasonInvalidToken           = PreconditionReason("invalid-token")
	PreconditionReasonTokenExpired           = PreconditionReason("token-expired")
	PreconditionReasonAccountIsLinked        = PreconditionReason("account-is-linked")
)

type PreconditionError struct {
	// The reason of PreconditionFailure, example: user-name-duplicated
	Reason      PreconditionReason
	Subject     string
	Description string
	Metadata    map[string]interface{}
}

func NewPreconditionError(reason PreconditionReason, subject string, description string, md map[string]interface{}) *PreconditionError {
	return &PreconditionError{
		Reason:      reason,
		Subject:     subject,
		Description: description,
		Metadata:    md,
	}
}

func (e *PreconditionError) Error() string {
	if e == nil {
		return "nil"
	}
	return fmt.Sprintf("PreconditionError: {type: %s, subject: %s, description: %s}", e.Reason, e.Subject, e.Description)
}

func (e *PreconditionError) Is(target error) bool {
	t, ok := target.(*PreconditionError)
	if !ok {
		return false
	}
	return e.Reason == t.Reason && e.Subject == t.Subject
}

func (e *PreconditionError) As(target interface{}) bool {
	t, ok := target.(*PreconditionError)
	if !ok {
		return false
	}
	*t = *e
	return true
}
