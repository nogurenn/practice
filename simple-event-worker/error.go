package simpleeventworker

import (
	"fmt"
)

var (
	ErrInputJSONIsNotArray = fmt.Errorf("input JSON is not an array")
)

type ErrUnsupportedEventType struct {
	Type string
}

func (e *ErrUnsupportedEventType) Error() string {
	return fmt.Sprintf(`unsupported event type: "%s"`, e.Type)
}

type ErrMissingFieldInEventPayloadField struct {
	Field string
}

func (e *ErrMissingFieldInEventPayloadField) Error() string {
	return fmt.Sprintf(`this event type requires the field: "%s"`, e.Field)
}

type ErrAccountDoesNotExist struct {
	AccountID string
}

func (e *ErrAccountDoesNotExist) Error() string {
	return fmt.Sprintf(`account with ID does not exist: "%s"`, e.AccountID)
}

type ErrAccountAlreadyExists struct {
	AccountID string
}

func (e *ErrAccountAlreadyExists) Error() string {
	return fmt.Sprintf(`account with ID already exists: "%s"`, e.AccountID)
}

type ErrCannotTransactWithRecalledAccount struct {
	AccountID string
}

func (e *ErrCannotTransactWithRecalledAccount) Error() string {
	return fmt.Sprintf(`cannot record charge or payment for recalled account with ID: "%s"`, e.AccountID)
}
