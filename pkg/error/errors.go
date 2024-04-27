package error

import (
	"fmt"
)

// RPError is an error that originates from your system (Responsible Party) that must not be communicated to the enduser as a BankID error.
type RPError struct {
	StatusCode int
	Message string
}

// BankIDError is an error originating from BankID's system that should be communicated to the enduser by the RP
type BankIDError struct {
	StatusCode int
	Message string
}

func (r RPError) Error() string {
	return fmt.Sprintf("statuscode: %d message: %s", r.StatusCode, r.Message)
}

func (r BankIDError) Error() string {
	return r.Message
}

const (
	RAF1  = "The user cancelled."
	RFA4  = "An identification or signing for this personal number is already started. Please try again."
	RFA5  = "Internal error. Please try again."
	RFA22 = "Unknown error. Please try again."
)

var (
	// RP must inform the user that an auth or sign order is already in progress for the user.
	// Message RFA4 should be used.
	ErrAlreadyInProgress = BankIDError{
		StatusCode: 208, // HTTP 208 - Already In Progress
		Message: RFA4,
	}

	// If an unknown errorCode is returned, RP should inform the user. Message RFA22 should be used.
	// RP should update their implementation to support the new errorCode as soon as possible.
	ErrUnknownErrorCode = BankIDError{
		StatusCode: 501, // HTTP 501 - Not Implemented
		Message: RFA22,
	}

	// RP must not automatically try again. This error may occur if the processing at RP or the communication is too slow.
	// RP must inform the user. Message RFA5 should be used.
	ErrRequestTimeout = BankIDError{
		StatusCode: 408, // HTTP 408 - Request Timeout
		Message: RFA5,
	}

	// RP must not automatically try again. RP must inform the user.
	// Message RFA5 should be used.
	ErrInternalError = BankIDError{
		StatusCode: 500, // HTTP 500 - Internal Server Error
		Message: RFA5,
	}

	// RP may try again without informing the user. If this error is returned repeatedly, RP must inform the user.
	// Message RFA5 should be used.
	ErrMaintenance = BankIDError{
		StatusCode: 503, // HTTP 503 - Service Unavailable
		Message: RFA5,
	}
)

var (
	// RP must not try the same request again. This is an internal error within the RP's system and must not be communicated to the user as a BankID error.
	ErrInvalidParameters = RPError{
		StatusCode: 400,  // HTTP 400 - Bad Request
		Message: "Invalid parameter. Invalid use of method. Potential causes include using an orderRef that previously resulted in a completed or failed order, orderRef that is too old, using the wrong certificate, oversized content, or non-JSON bodies. Internal error within the RP's system.",
	}

	// RP must not try the same request again. This is an internal error within the RP's system and must not be communicated to the user as a BankID error.
	ErrUnauthorized = RPError{
		StatusCode: 401, // HTTP 401 - Unauthorized
		Message: "RP does not have access to the service. Internal error within the RP's system.",
	}

	// RP must not try the same request again. This is an internal error within the RP's system and must not be communicated to the user as a BankID error.
	ErrNotFound = RPError{
		StatusCode: 404, // HTTP 404 - Not Found
		Message: "An erroneous URL path was used. Internal error within the RP's system.",
	}

	// RP must not try the same request again. This is an internal error within the RP's system and must not be communicated to the user as a BankID error.
	ErrMethodNotAllowed = RPError{
		StatusCode: 405, // HTTP 405 - Method Not Allowed
		Message: "Only http method POST is allowed. Internal error within the RP's system.",
	}

	// RP must not try the same request again. This is an internal error within the RP's system and must not be communicated to the user as a BankID error.
	ErrUnsupportedMediaType = RPError{
		StatusCode: 415,  // HTTP 415 - Unsupported Media Type
		Message: "Adding a 'charset' parameter after 'application/json' is not allowed. Internal error within the RP's system.",
	}
)
