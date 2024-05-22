package bankid

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// BankID is an interface for interacting with the BankID API.
// You can use it to authenticate users and sign using BankID.
// Documentation: https://www.bankid.com/en/utvecklare
type BankID interface {
	// 🗝️ Initiates an authentication order.
	// Use the collect method to query the status of the order. If the request is successful the response includes:
	// 	- orderRef
	// 	- autoStartToken
	// 	- qrStartToken
	// 	- qrStartSecret
	// Documentation: https://www.bankid.com/en/utvecklare/guider/teknisk-integrationsguide/graenssnittsbeskrivning/auth
	Auth(request AuthRequest) (*AuthResponse, error)

	// 🖋️ Initiates an signing order.
	// Use the collect method to query the status of the order. If the request is successful the response includes:
	// 	- orderRef
	// 	- autoStartToken
	// 	- qrStartToken
	// 	- qrStartSecret
	// Documentation: https://www.bankid.com/en/utvecklare/guider/teknisk-integrationsguide/graenssnittsbeskrivning/sign
	Sign(request SignRequest) (*SignResponse, error)

	// 🗝️ Initiates an authentication order when the user is talking to the RP over the phone.
	// Use the collect method to query the status of the order.
	// Documentation: https://www.bankid.com/en/utvecklare/guider/teknisk-integrationsguide/graenssnittsbeskrivning/phone-auth
	PhoneAuth(request PhoneAuthRequest) (*PhoneAuthResponse, error)

	// 🖋️ Initiates an signing order when the user is talking to the RP over the phone.
	// Use the collect method to query the status of the order.
	// Documentation: https://www.bankid.com/en/utvecklare/guider/teknisk-integrationsguide/graenssnittsbeskrivning/phone-sign
	PhoneSign(request PhoneSignRequest) (*PhoneSignResponse, error)

	// 🫳 Collects the result of a sign or auth order using orderRef as reference.
	// RP should keep on calling collect every two seconds if status is pending.
	// RP must abort if status indicates failed. The user identity is returned when complete.
	// Documentation: https://www.bankid.com/en/utvecklare/guider/teknisk-integrationsguide/graenssnittsbeskrivning/collect
	Collect(request CollectRequest) (*CollectResponse, error)

	// 🫳 Continuously calls the /collect endpoint (every 2 seconds) in a goroutine for as long as the order is pending
	// Collects the result of a sign or auth order using orderRef as reference
	// Will result in a succeeded or failed authentication. The user identity is returned when complete.
	//
	// Documentation: https://www.bankid.com/en/utvecklare/guider/teknisk-integrationsguide/graenssnittsbeskrivning/collect
	//
	// Example:
	// 		req := bankid.CollectRequest{
	// 			OrderRef: authResponse.OrderRef,
	// 		}
	//
	// 		go b.CollectRoutine(ctx, req, response)
	//
	// 		for {
	// 			select {
	// 			case collectResponse, ok := <-response:
	// 				// work with CollectResponse
	// 			case <-ctx.Done():
	// 			    return
	// 			}
	// 		}
	CollectRoutine(ctx context.Context, request CollectRequest, response chan *CollectResponse)

	// ✋ Cancels an ongoing sign or auth order.
	// This is typically used if the user cancels the order in your service or app.
	// Documentation: https://www.bankid.com/en/utvecklare/guider/teknisk-integrationsguide/graenssnittsbeskrivning/cancel
	Cancel(request CancelRequest) (*CancelResponse, error)
}

type bankid struct {
	config *RequestConfig
}

func New(input Parameters) (BankID, error) {
	err := input.Validate()
	if err != nil {
		return nil, fmt.Errorf("error validating parameters: %w", err)
	}

	c, err := newRequestConfig(input)
	if err != nil {
		return nil, fmt.Errorf("error creating new config: %w", err)
	}

	return &bankid{
		config: c,
	}, nil
}

// Generates the string that needs to be encoded into a QR code.
// The following pattern is used as a link in the QR code
//
//	`bankid.qrStartToken.time.qrAuthCode`
//
// BankID documentation: https://www.bankid.com/en/utvecklare/guider/teknisk-integrationsguide/qrkoder
func GenerateQrPayload(qrStartSecret string, qrStartToken string, timeInSeconds int) (string, error) {
	hash := hmac.New(sha256.New, []byte(qrStartSecret))
	_, err := hash.Write([]byte(fmt.Sprintf("%d", timeInSeconds)))
	if err != nil {
		return "", fmt.Errorf("error creating hash for qr: %w", err)
	}
	return fmt.Sprintf("bankid.%s.%d.%s", qrStartToken, timeInSeconds, hex.EncodeToString(hash.Sum(nil))), nil
}

// Initiates an authentication order. Use the collect method to query the status of the order.
func (b *bankid) Auth(req AuthRequest) (*AuthResponse, error) {	
	err := validate(
		validateRequired(req),
		validateEndUserIP(req.EndUserIP),
		validateRequirement(req.Requirement),
	)
	if err != nil {
		return nil, fmt.Errorf("input error: %w", err)
	}

	req, err = process[AuthRequest](
		processUserVisibleData(req.UserVisibleData),
		processUserNonVisibleData(req.UserNonVisibleData),
		processUserVisibleDataFormat(req.UserVisibleDataFormat),
	)
	if err != nil {
		return nil, fmt.Errorf("process error: %w", err)
	}

	return request[AuthResponse](RequestParameters{
		Path:   "/auth",
		Config: b.config,
		Body:   req,
	})
}

// Initiates an signing order. Use the collect method to query the status of the order.
func (b *bankid) Sign(req SignRequest) (*SignResponse, error) {
	err := validate(
		validateRequired(req),
		validateEndUserIP(req.EndUserIP),
		validateRequirement(req.Requirement),
	)
	if err != nil {
		return nil, fmt.Errorf("input error: %w", err)
	}

	req, err = process[SignRequest](
		processUserVisibleData(req.UserVisibleData),
		processUserNonVisibleData(req.UserNonVisibleData),
		processUserVisibleDataFormat(req.UserVisibleDataFormat),
	)
	if err != nil {
		return nil, fmt.Errorf("process error: %w", err)
	}

	return request[SignResponse](RequestParameters{
		Path:   "/sign",
		Config: b.config,
		Body:   req,
	})
}

// Initiates an authentication order when the user is talking to the RP over the phone.
func (b *bankid) PhoneAuth(req PhoneAuthRequest) (*PhoneAuthResponse, error) {
	err := validate(
		validateRequired(req),
		validatePersonalNumber(req.PersonalNumber),
		validateCallInitiator(req.CallInitiator),
		validateRequirement(req.Requirement),
	)
	if err != nil {
		return nil, fmt.Errorf("input error: %w", err)
	}

	req, err = process[PhoneAuthRequest](
		processUserVisibleData(req.UserVisibleData),
		processUserNonVisibleData(req.UserNonVisibleData),
		processUserVisibleDataFormat(req.UserVisibleDataFormat),
	)
	if err != nil {
		return nil, fmt.Errorf("process error: %w", err)
	}

	return request[PhoneAuthResponse](RequestParameters{
		Path:   "/phone/auth",
		Config: b.config,
		Body:   req,
	})
}

// Initiates an signing order when the user is talking to the RP over the phone.
func (b *bankid) PhoneSign(req PhoneSignRequest) (*PhoneSignResponse, error) {
	err := validate(
		validateRequired(req),
		validatePersonalNumber(req.PersonalNumber),
		validateCallInitiator(req.CallInitiator),
		validateRequirement(req.Requirement),
	)
	if err != nil {
		return nil, fmt.Errorf("input error: %w", err)
	}

	req, err = process[PhoneSignRequest](
		processUserVisibleData(req.UserVisibleData),
		processUserNonVisibleData(req.UserNonVisibleData),
		processUserVisibleDataFormat(req.UserVisibleDataFormat),
	)
	if err != nil {
		return nil, fmt.Errorf("process error: %w", err)
	}

	return request[PhoneSignResponse](RequestParameters{
		Path:   "/phone/sign",
		Config: b.config,
		Body:   req,
	})
}

// Cancels an ongoing sign or auth order.
func (b *bankid) Cancel(req CancelRequest) (*CancelResponse, error) {
	return request[CancelResponse](RequestParameters{
		Path:   "/cancel",
		Config: b.config,
		Body:   req,
	})
}

// Collects the result of a sign or auth order using orderRef as reference.
func (b *bankid) Collect(req CollectRequest) (*CollectResponse, error) {
	return request[CollectResponse](RequestParameters{
		Path:   "/collect",
		Config: b.config,
		Body:   req,
	})
}

// A goroutine that checks the /collect endpoint every 2 seconds and returns the response in a channel
// BankID reference: https://www.bankid.com/en/utvecklare/guider/teknisk-integrationsguide/graenssnittsbeskrivning/collect
func (b *bankid) CollectRoutine(ctx context.Context, request CollectRequest, response chan *CollectResponse) {
	defer close(response)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			collectResponse, err := b.Collect(request)
			if err != nil {
				fmt.Printf("Error collecting status: %v\n", err)
				return
			}

			response <- collectResponse

			if collectResponse.Status == Pending {
				time.Sleep(2 * time.Second)
				continue
			}

			return
		}
	}
}
