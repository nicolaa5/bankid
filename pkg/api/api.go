package api

import (
	"github.com/nicolaa5/bankid/internal/request"
	"github.com/nicolaa5/bankid/pkg/cfg"
)

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

type bankid struct {
	config cfg.Config
}

func New(config cfg.Config) (BankID, error) {

	request.New(config)

	return &bankid{
		config: config,
	}, nil
}


// Auth implements BankID.
func (b *bankid) Auth(personnummer string) error {
	panic("unimplemented")
}

// Cancel implements BankID.
func (b *bankid) Cancel(orderNummer string) error {
	panic("unimplemented")
}

// Collect implements BankID.
func (b *bankid) Collect(orderNummer string) error {
	panic("unimplemented")
}

// PhoneAuth implements BankID.
func (b *bankid) PhoneAuth() error {
	panic("unimplemented")
}

// PhoneSign implements BankID.
func (b *bankid) PhoneSign() error {
	panic("unimplemented")
}

// Sign implements BankID.
func (b *bankid) Sign() error {
	panic("unimplemented")
}