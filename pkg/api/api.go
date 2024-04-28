package api

import (
	"fmt"

	"github.com/nicolaa5/bankid/internal/client"
	"github.com/nicolaa5/bankid/pkg/parameters"
	"github.com/nicolaa5/bankid/pkg/request"
	"github.com/nicolaa5/bankid/pkg/response"
)

type BankID interface {
	// Initiates an authentication order.
	// Use the collect method to query the status of the order. If the request is successful the response includes:
	// 	- orderRef
	// 	- autoStartToken
	// 	- qrStartToken
	// 	- qrStartSecret
	// Documentation: https://www.bankid.com/en/utvecklare/guider/teknisk-integrationsguide/graenssnittsbeskrivning/auth
	Auth(request request.AuthRequest) (*response.AuthResponse, error)

	// Initiates an signing order.
	// Use the collect method to query the status of the order. If the request is successful the response includes:
	// 	- orderRef
	// 	- autoStartToken
	// 	- qrStartToken
	// 	- qrStartSecret
	// Documentation: https://www.bankid.com/en/utvecklare/guider/teknisk-integrationsguide/graenssnittsbeskrivning/sign
	Sign(request request.SignRequest) (*response.SignResponse, error)

	// Initiates an authentication order when the user is talking to the RP over the phone.
	// Use the collect method to query the status of the order.
	// Documentation: https://www.bankid.com/en/utvecklare/guider/teknisk-integrationsguide/graenssnittsbeskrivning/phone-auth
	PhoneAuth(request request.PhoneAuthRequest) (*response.PhoneAuthResponse, error)

	// Initiates an signing order when the user is talking to the RP over the phone.
	// Use the collect method to query the status of the order.
	// Documentation: https://www.bankid.com/en/utvecklare/guider/teknisk-integrationsguide/graenssnittsbeskrivning/phone-sign
	PhoneSign(request request.PhoneSignRequest) (*response.PhoneSignResponse, error)

	// Collects the result of a sign or auth order using orderRef as reference.
	// RP should keep on calling collect every two seconds if status is pending.
	// RP must abort if status indicates failed. The user identity is returned when complete.
	// Documentation: https://www.bankid.com/en/utvecklare/guider/teknisk-integrationsguide/graenssnittsbeskrivning/collect
	Collect(request request.CollectRequest) (*response.CollectResponse, error)

	// Cancels an ongoing sign or auth order.
	// This is typically used if the user cancels the order in your service or app.
	// Documentation: https://www.bankid.com/en/utvecklare/guider/teknisk-integrationsguide/graenssnittsbeskrivning/cancel
	Cancel(request request.CancelRequest) (*response.CancelResponse, error)
}

type bankid struct {
	client *client.Config
}

func New(config parameters.Parameters) (BankID, error) {
	err := config.Validate()
	if err != nil {
		return nil, fmt.Errorf("error validating parameters: %w", err)
	}

	c, err := client.New(config)
	if err != nil {
		return nil, fmt.Errorf("error initializing client: %w", err)
	}

	return &bankid{
		client: c,
	}, nil
}

// Initiates an authentication order. Use the collect method to query the status of the order.
func (b *bankid) Auth(request request.AuthRequest) (*response.AuthResponse, error) {
	return client.Request[response.AuthResponse](client.Parameters{
		Path:   "/auth",
		Config: b.client,
		Body:   request,
	})
}

// Initiates an signing order. Use the collect method to query the status of the order.
func (b *bankid) Sign(request request.SignRequest) (*response.SignResponse, error) {
	return client.Request[response.SignResponse](client.Parameters{
		Path:   "/sign",
		Config: b.client,
		Body:   request,
	})
}

// Initiates an authentication order when the user is talking to the RP over the phone.
func (b *bankid) PhoneAuth(request request.PhoneAuthRequest) (*response.PhoneAuthResponse, error) {
	return client.Request[response.PhoneAuthResponse](client.Parameters{
		Path:   "/cancel",
		Config: b.client,
		Body:   request,
	})
}

// Initiates an signing order when the user is talking to the RP over the phone.
func (b *bankid) PhoneSign(request request.PhoneSignRequest) (*response.PhoneSignResponse, error) {
	return client.Request[response.PhoneSignResponse](client.Parameters{
		Path:   "/phone/sign",
		Config: b.client,
		Body:   request,
	})
}

// Collects the result of a sign or auth order using orderRef as reference.
func (b *bankid) Collect(request request.CollectRequest) (*response.CollectResponse, error) {
	return client.Request[response.CollectResponse](client.Parameters{
		Path:   "/collect",
		Config: b.client,
		Body:   request,
	})
}

// Cancels an ongoing sign or auth order.
func (b *bankid) Cancel(request request.CancelRequest) (*response.CancelResponse, error) {
	return client.Request[response.CancelResponse](client.Parameters{
		Path:   "/cancel",
		Config: b.client,
		Body:   request,
	})
}
