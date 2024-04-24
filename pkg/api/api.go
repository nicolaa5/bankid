package api

type BankID interface {
	// auth
	Auth(personnummer string) error

	// sign
	Sign() error

	// phone/auth
	PhoneAuth() error

	// phone/sign
	PhoneSign() error

	// collect
	Collect() error

	// cancel
	Cancel() error
}