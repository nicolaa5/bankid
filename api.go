package bankid

import (
	"fmt"
)

type BankID interface {
	// Initiates an authentication order.
	// Use the collect method to query the status of the order. If the request is successful the response includes:
	// 	- orderRef
	// 	- autoStartToken
	// 	- qrStartToken
	// 	- qrStartSecret
	// Documentation: https://www.bankid.com/en/utvecklare/guider/teknisk-integrationsguide/graenssnittsbeskrivning/auth
	Auth(request AuthRequest) (*AuthResponse, error)

	// Initiates an signing order.
	// Use the collect method to query the status of the order. If the request is successful the response includes:
	// 	- orderRef
	// 	- autoStartToken
	// 	- qrStartToken
	// 	- qrStartSecret
	// Documentation: https://www.bankid.com/en/utvecklare/guider/teknisk-integrationsguide/graenssnittsbeskrivning/sign
	Sign(request SignRequest) (*SignResponse, error)

	// Initiates an authentication order when the user is talking to the RP over the phone.
	// Use the collect method to query the status of the order.
	// Documentation: https://www.bankid.com/en/utvecklare/guider/teknisk-integrationsguide/graenssnittsbeskrivning/phone-auth
	PhoneAuth(request PhoneAuthRequest) (*PhoneAuthResponse, error)

	// Initiates an signing order when the user is talking to the RP over the phone.
	// Use the collect method to query the status of the order.
	// Documentation: https://www.bankid.com/en/utvecklare/guider/teknisk-integrationsguide/graenssnittsbeskrivning/phone-sign
	PhoneSign(request PhoneSignRequest) (*PhoneSignResponse, error)

	// Collects the result of a sign or auth order using orderRef as reference.
	// RP should keep on calling collect every two seconds if status is pending.
	// RP must abort if status indicates failed. The user identity is returned when complete.
	// Documentation: https://www.bankid.com/en/utvecklare/guider/teknisk-integrationsguide/graenssnittsbeskrivning/collect
	Collect(request CollectRequest) (*CollectResponse, error)

	// Cancels an ongoing sign or auth order.
	// This is typically used if the user cancels the order in your service or app.
	// Documentation: https://www.bankid.com/en/utvecklare/guider/teknisk-integrationsguide/graenssnittsbeskrivning/cancel
	Cancel(request CancelRequest) (*CancelResponse, error)
}

type bankid struct {
	config *Config
}

func New(config Parameters) (BankID, error) {
	err := config.Validate()
	if err != nil {
		return nil, fmt.Errorf("error validating parameters: %w", err)
	}

	c, err := newConfig(config)
	if err != nil {
		return nil, fmt.Errorf("error creating new config: %w", err)
	}

	return &bankid{
		config: c,
	}, nil
}

// Initiates an authentication order. Use the collect method to query the status of the order.
func (b *bankid) Auth(req AuthRequest) (*AuthResponse, error) {
	return request[AuthResponse](ClientParameters{
		Path:   "/auth",
		Config: b.config,
		Body:   req,
	})
}

// Initiates an signing order. Use the collect method to query the status of the order.
func (b *bankid) Sign(req SignRequest) (*SignResponse, error) {
	return request[SignResponse](ClientParameters{
		Path:   "/sign",
		Config: b.config,
		Body:   req,
	})
}

// Initiates an authentication order when the user is talking to the RP over the phone.
func (b *bankid) PhoneAuth(req PhoneAuthRequest) (*PhoneAuthResponse, error) {
	return request[PhoneAuthResponse](ClientParameters{
		Path:   "/cancel",
		Config: b.config,
		Body:   req,
	})
}

// Initiates an signing order when the user is talking to the RP over the phone.
func (b *bankid) PhoneSign(req PhoneSignRequest) (*PhoneSignResponse, error) {
	return request[PhoneSignResponse](ClientParameters{
		Path:   "/phone/sign",
		Config: b.config,
		Body:   req,
	})
}

// Collects the result of a sign or auth order using orderRef as reference.
func (b *bankid) Collect(req CollectRequest) (*CollectResponse, error) {
	return request[CollectResponse](ClientParameters{
		Path:   "/collect",
		Config: b.config,
		Body:   req,
	})
}

// Cancels an ongoing sign or auth order.
func (b *bankid) Cancel(req CancelRequest) (*CancelResponse, error) {
	return request[CancelResponse](ClientParameters{
		Path:   "/cancel",
		Config: b.config,
		Body:   req,
	})
}
