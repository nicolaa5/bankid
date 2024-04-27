package api

import (
	"fmt"

	"github.com/nicolaa5/bankid/internal/client"
	"github.com/nicolaa5/bankid/pkg/cfg"
	"github.com/nicolaa5/bankid/pkg/request"
	"github.com/nicolaa5/bankid/pkg/response"
)

const (
	BankidURL     = "https://appapi2.bankid.com/rp/v6.0"
	BankidTestUrl = "https://appapi2.test.bankid.com/rp/v6.0"
)

type BankID interface {
	// Initiates an authentication order.
	// Use the collect method to query the status of the order. If the request is successful the response includes:
	// 	- orderRef
	// 	- autoStartToken
	// 	- qrStartToken
	// 	- qrStartSecret
	Auth(request request.AuthRequest) (response response.AuthResponse, err error)

	// Initiates an signing order.
	// Use the collect method to query the status of the order. If the request is successful the response includes:
	// 	- orderRef
	// 	- autoStartToken
	// 	- qrStartToken
	// 	- qrStartSecret
	Sign(request request.SignRequest) (response response.SignResponse, err error)

	// Initiates an authentication order when the user is talking to the RP over the phone.
	// Use the collect method to query the status of the order.
	PhoneAuth(request request.PhoneAuthRequest) (response response.PhoneAuthResponse, err error)

	// Initiates an authentication order when the user is talking to the RP over the phone.
	// Use the collect method to query the status of the order.
	PhoneSign(request request.PhoneSignRequest) (response response.PhoneSignResponse, err error)

	// Collects the result of a sign or auth order using orderRef as reference.
	// RP should keep on calling collect every two seconds if status is pending.
	// RP must abort if status indicates failed. The user identity is returned when complete.
	Collect(request request.CollectRequest) (response response.CollectResponse, err error)

	// Cancels an ongoing sign or auth order.
	// This is typically used if the user cancels the order in your service or app.
	Cancel(request request.CancelRequest) (response response.CancelResponse, err error)
}

type bankid struct {
	config cfg.Config
	client *client.Client
}

func New(config cfg.Config) (BankID, error) {
	c, err := client.New(config)
	if err != nil {
		return nil, fmt.Errorf("error initializing client: %w", err)
	}

	return &bankid{
		config: config,
		client: c,
	}, nil
}


// Auth implements BankID.
func (b *bankid) Auth(request request.AuthRequest) (response response.AuthResponse, err error) {
	b.client.Request()
}

// Cancel implements BankID.
func (*bankid) Cancel(request request.CancelRequest) (response response.CancelResponse, err error) {
	panic("unimplemented")
}

// Collect implements BankID.
func (*bankid) Collect(request request.CollectRequest) (response response.CollectResponse, err error) {
	panic("unimplemented")
}

// PhoneAuth implements BankID.
func (*bankid) PhoneAuth(request request.PhoneAuthRequest) (response response.PhoneAuthResponse, err error) {
	panic("unimplemented")
}

// PhoneSign implements BankID.
func (*bankid) PhoneSign(request request.PhoneSignRequest) (response response.PhoneSignResponse, err error) {
	panic("unimplemented")
}

// Sign implements BankID.
func (*bankid) Sign(request request.SignRequest) (response response.SignResponse, err error) {
	panic("unimplemented")
}
