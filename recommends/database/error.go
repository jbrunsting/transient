package database

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

type ConnectionError struct {
	InternalError string
}

func (e *ConnectionError) Error() string {
	return "Error connecting to database"
}

type NotFoundError struct {
	Object string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("Could not find %v", e.Object)
}

type DataViolation struct {
	Violation string
}

func (e *DataViolation) Error() string {
	return fmt.Sprintf("Data violation: %v", e.Violation)
}

type UniquenessViolation struct {
	Object string
}

func (e *UniquenessViolation) Error() string {
	return fmt.Sprintf("Uniqueness constraint violated for %v", e.Object)
}

type UnexpectedError struct {
	Action        string
	InternalError string
}

func (e *UnexpectedError) Error() string {
	return fmt.Sprintf("Got unexpected error when %v", e.Action)
}

func formatError(err error, object string, action string) error {
	if err == nil {
		return nil
	} else if err == sql.ErrNoRows {
		return &NotFoundError{Object: object}
	} else if sqlErr, ok := err.(*pq.Error); ok {
		switch sqlErr.Code.Class() {
		case "08":
			return &ConnectionError{InternalError: err.Error()}
		case "22":
			return &DataViolation{Violation: sqlErr.Detail}
		case "23":
			if sqlErr.Code.Name() == "unique_violation" {
				return &UniquenessViolation{Object: object}
			}
		}
	}

	return &UnexpectedError{Action: action, InternalError: err.Error()}
}
