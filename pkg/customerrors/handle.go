package customerrors

func AssignError(errorCode ErrorCode) BankIDError {
	switch errorCode {
	case AlreadyInProgress:
		return ErrAlreadyInProgress
	case RequestTimeout:
		return ErrRequestTimeout
	case InternalError:
		return ErrInternalError
	case Maintenance:
		return ErrMaintenance
	case InvalidParameters:
		return ErrInvalidParameters
	case Unauthorized:
		return ErrUnauthorized
	case NotFound:
		return ErrNotFound
	case MethodNotAllowed:
		return ErrMethodNotAllowed
	case UnsupportedMediaType:
		return ErrUnsupportedMediaType
	default:
		return ErrUnknownErrorCode
	}
}
