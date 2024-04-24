package api

// BankID is a go client for the BankID API
type BankID interface {
	// Initiates an authentication order. Use the collect method to query the status of the order. If the request is successful the response includes:
	// - orderRef
	// - autoStartToken
	// - qrStartToken
	// - qrStartSecret
	Auth(personnummer string) error

	// sign
	Sign() error

	// phone/auth
	PhoneAuth() error

	// phone/sign
	PhoneSign() error

	// collect
	Collect(orderNummer string) error

	// cancel
	Cancel(orderNummer string) error
}

func New(config Config) (*BankID, error) {
	if config.URL
}