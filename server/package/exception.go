package _package

import (
	"log"

	"github.com/gofiber/fiber/v3"
)

// ErrorType defines the type of error
type ErrorType string

const (
	BadRequestError     ErrorType = "BAD_REQUEST"
	ValidationError     ErrorType = "VALIDATION_ERROR"
	NotFoundError       ErrorType = "NOT_FOUND"
	UnauthorizedError   ErrorType = "UNAUTHORIZED"
	ForbiddenError      ErrorType = "FORBIDDEN"
	ConflictError       ErrorType = "CONFLICT"
	InternalServerError ErrorType = "INTERNAL_ERROR"
	ServiceError        ErrorType = "SERVICE_ERROR"
)

// AppException represents an application exception
type AppException struct {
	Type       ErrorType
	Message    string
	Details    string
	StatusCode int
}

// Error implements the error interface
func (e *AppException) Error() string {
	return e.Message
}

// NewAppException creates a new exception
func NewAppException(errorType ErrorType, message string, statusCode int, details string) *AppException {
	return &AppException{
		Type:       errorType,
		Message:    message,
		Details:    details,
		StatusCode: statusCode,
	}
}

// ThrowBadRequest throws a bad request exception
func ThrowBadRequest(message string, details string) *AppException {
	return NewAppException(BadRequestError, message, fiber.StatusBadRequest, details)
}

// ThrowValidationError throws a validation error exception
func ThrowValidationError(message string, details string) *AppException {
	return NewAppException(ValidationError, message, fiber.StatusBadRequest, details)
}

// ThrowNotFound throws a not found exception
func ThrowNotFound(message string) *AppException {
	return NewAppException(NotFoundError, message, fiber.StatusNotFound, "")
}

// ThrowUnauthorized throws an unauthorized exception
func ThrowUnauthorized(message string) *AppException {
	return NewAppException(UnauthorizedError, message, fiber.StatusUnauthorized, "")
}

// ThrowForbidden throws a forbidden exception
func ThrowForbidden(message string) *AppException {
	return NewAppException(ForbiddenError, message, fiber.StatusForbidden, "")
}

// ThrowConflict throws a conflict exception
func ThrowConflict(message string, details string) *AppException {
	return NewAppException(ConflictError, message, fiber.StatusConflict, details)
}

// ThrowInternalError throws an internal error exception
func ThrowInternalError(message string, details string) *AppException {
	log.Printf("[ERROR] Internal Error: %s - %s", message, details)
	return NewAppException(InternalServerError, message, fiber.StatusInternalServerError, details)
}

// ThrowServiceError throws a service error exception
func ThrowServiceError(message string, details string) *AppException {
	log.Printf("[ERROR] Service Error: %s - %s", message, details)
	return NewAppException(ServiceError, message, fiber.StatusServiceUnavailable, details)
}

// HandleException handles an exception and sends the response
func HandleException(c fiber.Ctx, err error) error {
	// Check if the error is an AppException
	if appErr, ok := err.(*AppException); ok {
		response := ErrorResponse(appErr.Message, string(appErr.Type), appErr.Details)
		return SendJSON(c, appErr.StatusCode, response)
	}

	// If not an AppException, treat it as an internal error
	log.Printf("[ERROR] Unexpected error: %v", err)
	response := ErrorResponse("Internal server error", string(InternalServerError), err.Error())
	return SendJSON(c, fiber.StatusInternalServerError, response)
}

// ExceptionMiddleware middleware for handling exceptions
func ExceptionMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		err := c.Next()

		if err != nil {
			return HandleException(c, err)
		}

		return nil
	}
}

// PanicRecoveryMiddleware middleware for catching panics
func PanicRecoveryMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[PANIC] Panic occurred: %v", r)
				
				response := ErrorResponse("Internal server error", string(InternalServerError), "Panic occurred")
				_ = SendJSON(c, fiber.StatusInternalServerError, response)
			}
		}()

		return c.Next()
	}
}

// Success helper functions for sending success responses

// SendSuccess sends a success response (HTTP 200)
func SendSuccess(c fiber.Ctx, message string, data interface{}) error {
	response := SuccessResponse(message, data)
	return SendJSON(c, fiber.StatusOK, response)
}

// SendCreated sends a created response (HTTP 201)
func SendCreated(c fiber.Ctx, message string, data interface{}) error {
	response := SuccessResponse(message, data)
	return SendJSON(c, fiber.StatusCreated, response)
}
