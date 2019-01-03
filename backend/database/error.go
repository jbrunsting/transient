package database

import (
	"fmt"
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

type UnexpectedError struct {
	Action string
	InternalError string
}

func (e *UnexpectedError) Error() string {
	return fmt.Sprintf("Got unexpected error when %v", e.Action)
}
