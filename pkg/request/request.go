package request

// RequestBody is an interface for all BankID requests.
type RequestBody interface {
	// Marshal returns the JSON encoded body of the request.
	Marshal() ([]byte, error)
}
