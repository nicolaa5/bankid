package customerrors

import (
	"fmt"
)

// BankIDError is an error returned by BankID that should be communicated to the enduser, or handled by the RP.
type BankIDError struct {
	StatusCode int `json:"statusCode"`
	Details    string `json:"details"`
	ErrorCode  ErrorCode `json:"errorCode"`
}

func (r BankIDError) Error() string {
	return fmt.Sprintf("errorcode: %s statuscode: %d details: %s", r.ErrorCode, r.StatusCode, r.Details)
}

const (
	RAF1  = "The user cancelled."
	RFA4  = "An identification or signing for this personal number is already started. Please try again."
	RFA5  = "Internal error. Please try again."
	RFA22 = "Unknown error. Please try again."
)

type ErrorCode string

const (
	AlreadyInProgress    ErrorCode = "alreadyInProgress"
	UnknownErrorCode     ErrorCode = "unknownErrorCode"
	RequestTimeout       ErrorCode = "requestTimeout"
	InternalError        ErrorCode = "internalError"
	Maintenance          ErrorCode = "maintenance"
	InvalidParameters    ErrorCode = "invalidParameters"
	Unauthorized         ErrorCode = "unauthorized"
	NotFound             ErrorCode = "notFound"
	MethodNotAllowed     ErrorCode = "methodNotAllowed"
	UnsupportedMediaType ErrorCode = "unsupportedMediaType"
)

var (
	// RP must inform the user that an auth or sign order is already in progress for the user.
	// Message RFA4 should be used.
	ErrAlreadyInProgress = BankIDError{
		StatusCode: 208, // HTTP 208 - Already In Progress
		Details:    RFA4,
		ErrorCode: AlreadyInProgress,
	}

	// If an unknown errorCode is returned, RP should inform the user. Message RFA22 should be used.
	// RP should update their implementation to support the new errorCode as soon as possible.
	ErrUnknownErrorCode = BankIDError{
		StatusCode: 501, // HTTP 501 - Not Implemented
		Details:    RFA22,
		ErrorCode:  UnknownErrorCode,
	}

	// RP must not automatically try again. This error may occur if the processing at RP or the communication is too slow.
	// RP must inform the user. Message RFA5 should be used.
	ErrRequestTimeout = BankIDError{
		StatusCode: 408, // HTTP 408 - Request Timeout
		Details:    RFA5,
		ErrorCode:  RequestTimeout,
	}

	// RP must not automatically try again. RP must inform the user.
	// Message RFA5 should be used.
	ErrInternalError = BankIDError{
		StatusCode: 500, // HTTP 500 - Internal Server Error
		Details:    RFA5,
		ErrorCode:  InternalError,
	}

	// RP may try again without informing the user. If this error is returned repeatedly, RP must inform the user.
	// Message RFA5 should be used.
	ErrMaintenance = BankIDError{
		StatusCode: 503, // HTTP 503 - Service Unavailable
		Details:    RFA5,
		ErrorCode:  Maintenance,
	}
	
	// RP must not try the same request again. This is an internal error within the RP's system and must not be communicated to the user as a BankID error.
	ErrInvalidParameters = BankIDError{
		StatusCode: 400, // HTTP 400 - Bad Request
		Details:    "Invalid parameter. Invalid use of method. Potential causes include using an orderRef that previously resulted in a completed or failed order, orderRef that is too old, using the wrong certificate, oversized content, or non-JSON bodies. Internal error within the RP's system.",
		ErrorCode:  InvalidParameters,
	}

	// RP must not try the same request again. This is an internal error within the RP's system and must not be communicated to the user as a BankID error.
	ErrUnauthorized = BankIDError{
		StatusCode: 401, // HTTP 401 - Unauthorized
		Details:    "RP does not have access to the service. Internal error within the RP's system.",
		ErrorCode:  Unauthorized,
	}

	// RP must not try the same request again. This is an internal error within the RP's system and must not be communicated to the user as a BankID error.
	ErrNotFound = BankIDError{
		StatusCode: 404, // HTTP 404 - Not Found
		Details:    "An erroneous URL path was used. Internal error within the RP's system.",
		ErrorCode:  NotFound,
	}

	// RP must not try the same request again. This is an internal error within the RP's system and must not be communicated to the user as a BankID error.
	ErrMethodNotAllowed = BankIDError{
		StatusCode: 405, // HTTP 405 - Method Not Allowed
		Details:    "Only http method POST is allowed. Internal error within the RP's system.",
		ErrorCode:  MethodNotAllowed,
	}

	// RP must not try the same request again. This is an internal error within the RP's system and must not be communicated to the user as a BankID error.
	ErrUnsupportedMediaType = BankIDError{
		StatusCode: 415, // HTTP 415 - Unsupported Media Type
		Details:    "Adding a 'charset' parameter after 'application/json' is not allowed. Internal error within the RP's system.",
		ErrorCode:  UnsupportedMediaType,
	}
)
