package errs

import "fmt"

const FAIL = "[FAIL]"

type BadRequestError struct {
    Code    int
    Message string
}

func (e *BadRequestError) Error() string {
    return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func NewBadRequestError(message string) *BadRequestError {
    return &BadRequestError{
        Code:    400,
        Message: message,
    }
}

type NotFoundError struct {
    Code    int
    Message string
}

func (e *NotFoundError) Error() string {
    return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func NewNotFoundError(message string) *NotFoundError {
    return &NotFoundError{
        Code:    404,
        Message: message,
    }
}

type ValidationError struct {
    Code    int
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func NewValidationError(message string) *ValidationError {
    return &ValidationError{
        Code:    400,
        Message: message,
    }
}

type DuplicateEntryError struct {
    Code    int
    Message string
}

func (e *DuplicateEntryError) Error() string {
    return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func NewDuplicateEntryError(message string) *DuplicateEntryError {
    return &DuplicateEntryError{
        Code:    409,
        Message: message,
    }
}

type ForeignKeyViolationError struct {
    Code    int
    Message string
}

func (e *ForeignKeyViolationError) Error() string {
    return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func NewForeignKeyViolationError(message string) *ForeignKeyViolationError {
    return &ForeignKeyViolationError{
        Code:    400,
        Message: message,
    }
}

type DatabaseConnectionError struct {
    Code    int
    Message string
}

func (e *DatabaseConnectionError) Error() string {
    return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func NewDatabaseConnectionError(message string) *DatabaseConnectionError {
    return &DatabaseConnectionError{
        Code:    503,
        Message: message,
    }
}

type DatabaseError struct {
    Code    int
    Message string
}

func (e *DatabaseError) Error() string {
    return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func NewDatabaseError(message string) *DatabaseError {
    return &DatabaseError{
        Code:    500,
        Message: message,
    }
}

type InternalServerError struct {
    Code    int
    Message string
}

func (e *InternalServerError) Error() string {
    return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func NewInternalServerError(message string) *InternalServerError {
    return &InternalServerError{
        Code: 500,
        Message: message,
    }
}

type ForbiddenError struct {
    Code    int
    Message string
}

func (e *ForbiddenError) Error() string {
    return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func NewForbiddenError(message string) *ForbiddenError {
    return &ForbiddenError{
        Code: 403,
        Message: message,
    }
}


type UnauthorizedError struct {
    Code    int
    Message string
}

func (e *UnauthorizedError) Error() string {
    return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func NewUnauthorizedError(message string) *UnauthorizedError {
    return &UnauthorizedError{
        Code: 401,
        Message: message,
    }
}

type EntityTooLargeError struct {
    Code    int
    Message string
}

func (e *EntityTooLargeError) Error() string {
    return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func NewEntityTooLargeError(message string) *EntityTooLargeError {
    return &EntityTooLargeError{
        Code: 413,
        Message: message,
    }
}
