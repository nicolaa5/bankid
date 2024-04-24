package api

import (
	"fmt"

	"github.com/nicolaa5/bankid/internal/req"
	"github.com/nicolaa5/bankid/pkg/cfg"
)

const bankIDURL = "https://appapi2.bankid.com/rp/v6.0"

type BankID interface {
	// Initiates an authentication order. Use the collect method to query the status of the order. If the request is successful the response includes:
	// 	- orderRef
	// 	- autoStartToken
	// 	- qrStartToken
	// 	- qrStartSecret
	Auth(endUserIPAdress string) error

	// Initiates an signing order. Use the collect method to query the status of the order. If the request is successful the response includes:
	// 	- orderRef
	// 	- autoStartToken
	// 	- qrStartToken
	// 	- qrStartSecret
	Sign(endUserIPAdress string) error

	// phone/auth
	PhoneAuth(personnummer string) error

	// phone/sign
	PhoneSign(personnummer string) error

	// collect
	Collect(orderNummer string) error

	// cancel
	Cancel(orderNummer string) error
}

type bankid struct {
	config  cfg.Config
	request req.Client
}

// Auth implements BankID.
func (*bankid) Auth(endUserIPAdress string) error {
	panic("unimplemented")
}

// Cancel implements BankID.
func (*bankid) Cancel(orderNummer string) error {
	panic("unimplemented")
}

// Collect implements BankID.
func (*bankid) Collect(orderNummer string) error {
	panic("unimplemented")
}

// PhoneAuth implements BankID.
func (*bankid) PhoneAuth(personnummer string) error {
	panic("unimplemented")
}

// PhoneSign implements BankID.
func (*bankid) PhoneSign(personnummer string) error {
	panic("unimplemented")
}

// Sign implements BankID.
func (*bankid) Sign(endUserIPAdress string) error {
	panic("unimplemented")
}

func New(config cfg.Config) (BankID, error) {

	r, err := req.New(config)
	if err != nil {
		return nil, fmt.Errorf("error initializing: %w", err)
	}

	return &bankid{
		config: config,
	}, nil
}
