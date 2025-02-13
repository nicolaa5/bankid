package bankid

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

const (
	BankIDURL            = "https://appapi2.bankid.com/rp/v6.0"
	BankIDTestUrl        = "https://appapi2.test.bankid.com/rp/v6.0"
	BankIDTestPassphrase = "qwerty123"
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
	//
	// Documentation: https://www.bankid.com/en/utvecklare/guider/teknisk-integrationsguide/graenssnittsbeskrivning/auth
	//
	// Request flow:
	// 	- Starting the request: `hintCode: outstandingTransaction`
	// 	- User needs to provide the pin to authenticate themselves: `hintCode: userSign`
	// 	- User has authenticated themselves successfully: `status: complete`
	Auth(ctx context.Context, request AuthRequest) (*AuthResponse, error)

	// 🖋️ Initiates an signing order.
	// Use the collect method to query the status of the order. If the request is successful the response includes:
	// 	- orderRef
	// 	- autoStartToken
	// 	- qrStartToken
	// 	- qrStartSecret
	//
	// Documentation: https://www.bankid.com/en/utvecklare/guider/teknisk-integrationsguide/graenssnittsbeskrivning/sign
	//
	// Request flow:
	// 	- Starting the request: `hintCode: outstandingTransaction`
	// 	- User needs to provide the pin to sign the document: `hintCode: userSign`
	// 	- User has signed the document successfully: `status: complete`
	Sign(ctx context.Context, request SignRequest) (*SignResponse, error)

	// 🗝️ Initiates an authentication order when the user is talking to the RP over the phone.
	// Use the collect method to query the status of the order.
	//
	// Documentation: https://www.bankid.com/en/utvecklare/guider/teknisk-integrationsguide/graenssnittsbeskrivning/phone-auth
	//
	// Request flow:
	// 	- Starting the request: `hintCode: outstandingTransaction`
	// 	- User needs to confirm that they called, or were called by RP: `hintCode: userCallConfirm`
	// 	- User needs to provide the pin to authenticate themselves: `hintCode: userSign`
	// 	- User has authenticated themselves successfully: `status: complete`
	PhoneAuth(ctx context.Context, request PhoneAuthRequest) (*PhoneAuthResponse, error)

	// 🖋️ Initiates an signing order when the user is talking to the RP over the phone.
	// Use the collect method to query the status of the order.
	//
	// Documentation: https://www.bankid.com/en/utvecklare/guider/teknisk-integrationsguide/graenssnittsbeskrivning/phone-sign
	//
	// Request flow:
	// 	- Starting the request: `hintCode: outstandingTransaction`
	// 	- User needs to confirm that they called, or were called by RP: `hintCode: userCallConfirm`
	// 	- User needs to provide the pin to sign the document: `hintCode: userSign`
	// 	- User has signed the document successfully: `status: complete`
	PhoneSign(ctx context.Context, request PhoneSignRequest) (*PhoneSignResponse, error)

	// 🫳 Collects the result of a sign or auth order using orderRef as reference.
	// RP should keep on calling collect every two seconds if status is pending.
	// RP must abort if status indicates failed. The user identity is returned when complete.
	//
	// Documentation: https://www.bankid.com/en/utvecklare/guider/teknisk-integrationsguide/graenssnittsbeskrivning/collect
	Collect(ctx context.Context, request CollectRequest) (*CollectResponse, error)

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
	//
	// Documentation: https://www.bankid.com/en/utvecklare/guider/teknisk-integrationsguide/graenssnittsbeskrivning/cancel
	Cancel(ctx context.Context, request CancelRequest) (*CancelResponse, error)
}

type bankid struct {
	config *RequestConfig
}

func New(config Config) (BankID, error) {
	err := config.Validate()
	if err != nil {
		return nil, fmt.Errorf("error validating input config: %w", err)
	}

	config.UseDefault()

	c, err := newRequestConfig(config)
	if err != nil {
		return nil, fmt.Errorf("error creating new config: %w", err)
	}

	return &bankid{
		config: c,
	}, nil
}

// Returns a default Test BankID interface with SSL/CA certificates and password
func NewTestDefault() (BankID, error) {
	config := Config{
		URL: BankIDTestUrl,
		Certificate: P12Cert{
			Passphrase:    BankIDTestPassphrase,
			Certificate:   P12TestCertificate,
			CACertificate: CATestCertificate,
		},
	}

	c, err := newRequestConfig(config)
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
func (b *bankid) Auth(ctx context.Context, req AuthRequest) (*AuthResponse, error) {
	err := validate(
		validateRequired(req),
		validateEndUserIP(req.EndUserIP),
		validateRequirement(req.Requirement),
	)
	if err != nil {
		return nil, err
	}

	req, err = process[AuthRequest](req,
		processUserVisibleData(req.UserVisibleData),
		processUserNonVisibleData(req.UserNonVisibleData),
		processUserVisibleDataFormat(req.UserVisibleDataFormat),
	)
	if err != nil {
		return nil, fmt.Errorf("process error: %w", err)
	}

	return request[AuthResponse](ctx, RequestParameters{
		Path:   "/auth",
		Config: b.config,
		Body:   req,
	})
}

// Initiates an signing order. Use the collect method to query the status of the order.
func (b *bankid) Sign(ctx context.Context, req SignRequest) (*SignResponse, error) {
	err := validate(
		validateRequired(req),
		validateEndUserIP(req.EndUserIP),
		validateRequirement(req.Requirement),
	)
	if err != nil {
		return nil, err
	}

	req, err = process[SignRequest](req,
		processUserVisibleData(req.UserVisibleData),
		processUserNonVisibleData(req.UserNonVisibleData),
		processUserVisibleDataFormat(req.UserVisibleDataFormat),
	)
	if err != nil {
		return nil, err
	}

	return request[SignResponse](ctx, RequestParameters{
		Path:   "/sign",
		Config: b.config,
		Body:   req,
	})
}

// Initiates an authentication order when the user is talking to the RP over the phone.
func (b *bankid) PhoneAuth(ctx context.Context, req PhoneAuthRequest) (*PhoneAuthResponse, error) {
	err := validate(
		validateRequired(req),
		validatePersonalNumber(req.PersonalNumber),
		validateCallInitiator(req.CallInitiator),
		validateRequirement(req.Requirement),
	)
	if err != nil {
		return nil, err
	}

	req, err = process[PhoneAuthRequest](req,
		processUserVisibleData(req.UserVisibleData),
		processUserNonVisibleData(req.UserNonVisibleData),
		processUserVisibleDataFormat(req.UserVisibleDataFormat),
	)
	if err != nil {
		return nil, fmt.Errorf("process error: %w", err)
	}

	return request[PhoneAuthResponse](ctx, RequestParameters{
		Path:   "/phone/auth",
		Config: b.config,
		Body:   req,
	})
}

// Initiates an signing order when the user is talking to the RP over the phone.
func (b *bankid) PhoneSign(ctx context.Context, req PhoneSignRequest) (*PhoneSignResponse, error) {
	err := validate(
		validateRequired(req),
		validatePersonalNumber(req.PersonalNumber),
		validateCallInitiator(req.CallInitiator),
		validateRequirement(req.Requirement),
	)
	if err != nil {
		return nil, err
	}

	req, err = process[PhoneSignRequest](req,
		processUserVisibleData(req.UserVisibleData),
		processUserNonVisibleData(req.UserNonVisibleData),
		processUserVisibleDataFormat(req.UserVisibleDataFormat),
	)
	if err != nil {
		return nil, fmt.Errorf("process error: %w", err)
	}

	return request[PhoneSignResponse](ctx, RequestParameters{
		Path:   "/phone/sign",
		Config: b.config,
		Body:   req,
	})
}

// Cancels an ongoing sign or auth order.
func (b *bankid) Cancel(ctx context.Context, req CancelRequest) (*CancelResponse, error) {
	return request[CancelResponse](ctx, RequestParameters{
		Path:   "/cancel",
		Config: b.config,
		Body:   req,
	})
}

// Collects the result of a sign or auth order using orderRef as reference.
func (b *bankid) Collect(ctx context.Context, req CollectRequest) (*CollectResponse, error) {
	return request[CollectResponse](ctx, RequestParameters{
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
			collectResponse, err := b.Collect(ctx, request)
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
