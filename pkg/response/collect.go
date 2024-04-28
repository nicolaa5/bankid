package response

import "encoding/json"

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
