package bankid

import "encoding/json"

// ResponseBody is an interface for all successfull BankID responses.
type ResponseBody interface {
	// Unmarshal parses the JSON-encoded data and stores the result in the value pointed to by v.
	Unmarshal(data []byte) error
}

type ErrorResponseBody struct {
	ErrorCode int    `json:"errorCode"`
	Details   string `json:"details"`
}


// Response received from the auth endpoint, example of the response body: 
// {
//     "orderRef": "4820e3da-fbd7-45c0-aa1c-9d28d308c63b",
//     "autoStartToken": "f5071e97-ad0c-45ff-bc05-ca4bd8f14f84",
//     "qrStartToken": "a9002853-0445-4021-be15-9e373b71634a",
//     "qrStartSecret": "1238c8af-66d1-4c4a-8c00-b77deabeea98"
// }
type AuthResponse struct {
	// Used to collect the status of the order.
	OrderRef string `json:"orderRef"`

	// Used to compile the start url according to launching.
	// See https://www.bankid.com/utvecklare/guider/teknisk-integrationsguide/programstart
	AutoStartToken string `json:"autoStartToken"`

	// Used to compute the animated QR code.
	QrStartToken string `json:"qrStartToken"`

	// Used to compute the animated QR code.
	QrStartSecret string `json:"qrStartSecret"`
}

func (r AuthResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &r)
}

type SignResponse struct {
	// Used to collect the status of the order.
	OrderRef string `json:"orderRef"`
}

func (r SignResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &r)
}

type PhoneAuthResponse struct {
	// Used to collect the status of the order.
	OrderRef string `json:"orderRef"`
}

func (r PhoneAuthResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &r)
}

type PhoneSignResponse struct {
	// Used to collect the status of the order.
	OrderRef string `json:"orderRef"`
}

func (r PhoneSignResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &r)
}

// Response received from the collect endpoint, example of the collect response body:
// {
//     "orderRef": "b7e32e4f-0c1f-472b-9111-40dc856464b4",
//     "status": "pending",
//     "hintCode": "outstandingTransaction"
// }
type CollectResponse struct {
	OrderRef       string         `json:"orderRef"`
	Status         Status         `json:"status"`
	HintCode       HintCode       `json:"hintCode,omitempty"`
	CompletionData CompletionData `json:"completionData,omitempty"`
}

func (r CollectResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &r)
}

type User struct {
	PersonalNumber string `json:"personalNumber"`
	Name           string `json:"name"`
	GivenName      string `json:"givenName"`
	Surname        string `json:"surname"`
}

type Device struct {
	IpAddress string `json:"ipAddress"`
	Uhi       string `json:"uhi"`
}

type CompletionData struct {
	User            User   `json:"user"`
	Device          Device `json:"device"`
	BankIdIssueDate string `json:"bankIdIssueDate"`
	StepUp          bool   `json:"stepUp"`
	Signature       string `json:"signature"`
	OcspResponse    string `json:"ocspResponse"`
}

type Status string

const (
	Pending  Status = "pending"
	Failed   Status = "failed"
	Complete Status = "complete"
)

type HintCode string

const (
	// Order is pending. The BankID app has not yet received the order. The hintCode will later change to noClient, started or userSign.
	OutstandingTransaction HintCode = "outstandingTransaction"

	// Order is pending. The client has not yet received the order.
	NoClient HintCode = "noClient"

	// Order is pending. A BankID client has launched with autostarttoken but a usable ID has not yet been found in the client.
	// When the client launches there may be a short delay until all IDs are registered. The user may not have any usable IDs, or is yet to insert their smart card.
	Started HintCode = "started"

	// Order is pending. A client has launched and received the order but additional steps for providing MRTD information is required to proceed with the order.
	UserMrtd HintCode = "userMrtd"

	// Order is waiting for the user to confirm that they have received this order while in a call with the RP.
	UserCallConfirm HintCode = "userCallConfirm"

	// Order is pending. The BankID client has received the order.
	UserSign HintCode = "userSign"

	// The order has expired. The BankID security app/program did not launch, the user did not finalize the signing or the RP called collect too late.
	ExpiredTransaction HintCode = "expiredTransaction"

	// This error is returned if:
	// 	1. The user has entered the wrong PIN code too many times. The BankID cannot be used.
	// 	2. The user’s BankID is blocked.
	// 	3. The user’s BankID is invalid.
	CertificateErr HintCode = "certificateErr"

	// The order was cancelled by the user. userCancel may also be returned in some rare cases related to other user interactions.
	UserCancel HintCode = "userCancel"

	// The order was cancelled. The system received a new order for the user.
	Cancelled HintCode = "cancelled"

	// The user did not provide their ID or the client did not launch within a certain time limit. Potential
	// causes are:
	// 	1. RP did not use autoStartToken when launching the BankID security app. RP must correct this in their implementation.
	// 	2. Client software was not installed or other problem with the user’s device.
	StartFailed HintCode = "startFailed"
)

// A successful response contains an empty JSON object.
type CancelResponse struct{}

func (r CancelResponse) Unmarshal(data []byte) error {
	return nil
}
