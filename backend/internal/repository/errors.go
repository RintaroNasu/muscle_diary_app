package repository

import "errors"

var ErrFKViolation = errors.New("foreign key violation")
var ErrNotFound = errors.New("record not found")

type ConstraintError struct {
	Constraint string
}

func (e *ConstraintError) Error() string {
	return "constraint violation: " + e.Constraint
}
