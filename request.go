package bankid

import (
	"encoding/json"
)

// RequestBody is an interface for all BankID requests.
type RequestBody interface {
	// Marshal returns the JSON encoded body of the request.
	Marshal() ([]byte, error)
}

type AuthRequest struct {
	// Required: The user IP address as seen by RP. String. IPv4 and IPv6 is allowed.
	// Correct IP address must be the IP address representing the user agent (the end user device) as seen by the RP.
	// In case of inbound proxy, special considerations may need to be taken into account to get the correct address.
	// In some use cases the IP address is not available, for instance in voice-based services.
	// In these cases, the internal representation of those systems’ IP address may be used.
	EndUserIP string `json:"endUserIp"`

	// Optional: Requirements on how the auth order must be performed.
	Requirement Requirement `json:"requirement,omitempty"`

	// Optional: Text displayed to the user during authentication with BankID, with the purpose of providing context for the authentication
	// and to enable users to detect identification errors and averting fraud attempts.
	// The text can be formatted using CR, LF and CRLF for new lines. The text must be encoded as UTF-8 and then base 64 encoded. 1—1 500 characters after base 64 encoding.
	UserVisibleData string `json:"userVisibleData,omitempty"`

	// Optional: Data is not displayed to the user. String. The value must be base 64-encoded. 1-1 500 characters after base 64-encoding.
	UserNonVisibleData string `json:"userNonVisibleData,omitempty"`

	// Optional: If present, and set to “simpleMarkdownV1”, this parameter indicates that userVisibleData holds formatting characters which,
	// will potentially make the text displayed to the user nicer to look at.
	// For instructions check out https://www.bankid.com/utvecklare/guider/formatera-text
	UserVisibleDataFormat string `json:"userVisibleDataFormat,omitempty"`
	
	// Optional: Orders started on the same device (started with autostart token) will call this URL when the order is completed, 
	// ignoring any return URL provided in the start URL when the BankID app was launched.
	// If the user has an old version of the  BankID mobile client that doesn’t support getting the returnUrl 
	// from the server the user will be informed to update their client and the order will be cancelled.
	// The return URL you provide should include a nonce to the session. 
	// When the user returns to your app/webpage, your service should verify that the BankID transaction was completed successfully 
	// and that the device receiving the returnUrl is the same as the device that started the transaction.
	// Using this parameter will make your service more secure and strengthen the channel binding between you and the user.
	// Check that cookie/ip hasn’t changed from starting page to returnUrl page. String 1- 512 character.
	ReturnUrl string `json:"returnUrl,omitempty"`
	
	// Optional: If this is set to true, a risk indication will be included in the collect response when the order completes.
	// The risk indication requires that the endUserIp is correct. 
	// If your service provides an incorrect IP-address, legitimate transaction might be blocked and/or you might receive incorrect risk-data in the response.
	// Boolean. Default is false.
	ReturnRisk string `json:"returnRisk,omitempty"`
}

func (r AuthRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type SignRequest struct {
	// Required: The user IP address as seen by RP. String. IPv4 and IPv6 is allowed.
	// Correct IP address must be the IP address representing the user agent (the end user device) as seen by the RP. In case of inbound proxy, special considerations may need to be taken into account to get the correct address.
	// In some use cases the IP address is not available, for instance in voice-based services. In these cases, the internal representation of those systems’ IP address may be used.
	EndUserIP string `json:"endUserIp"`

	// Required: Text to be displayed to the user. String. The text can be formatted using CR, LF and CRLF for new lines.
	// The text must be encoded as UTF-8 and then base 64 encoded. 1 – 40,000 characters after base 64 encoding.
	UserVisibleData string `json:"userVisibleData"`

	// Optional: Requirements on how the auth order must be performed. See section Requirements below for more details.
	Requirement Requirement `json:"requirement,omitempty"`

	// Optional: Data is not displayed to the user. String. The value must be base 64-encoded. 1-1 500 characters after base 64-encoding.
	UserNonVisibleData string `json:"userNonVisibleData,omitempty"`

	// Optional: If present, and set to “simpleMarkdownV1”, this parameter indicates that userVisibleData holds formatting characters which,
	// will potentially make the text displayed to the user nicer to look at.
	// For instructions check out https://www.bankid.com/utvecklare/guider/formatera-text
	UserVisibleDataFormat string `json:"userVisibleDataFormat,omitempty"`
}

func (r SignRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Initiates an authentication order when the user is talking to the RP over the phone.
// Use the collect method to query the status of the order.
type PhoneAuthRequest struct {
	// Required: The personal number of the user. String. 12 digits.
	PersonalNumber string `json:"personalNumber"`

	// Required: Indicate if the user or the RP initiated the phone call.
	// 	- user: user called the RP
	// 	- RP: RP called the user
	CallInitiator string `json:"callInitiator"`

	// Optional: Requirements on how the auth order must be performed.
	Requirement Requirement `json:"requirement,omitempty"`

	// Optional: Text displayed to the user during authentication with BankID, with the purpose of providing context for the authentication
	// and to enable users to detect identification errors and averting fraud attempts.
	// The text can be formatted using CR, LF and CRLF for new lines. The text must be encoded as UTF-8 and then base 64 encoded. 1—1 500 characters after base 64 encoding.
	UserVisibleData string `json:"userVisibleData,omitempty"`

	// Optional: Data is not displayed to the user. String. The value must be base 64-encoded. 1-1 500 characters after base 64-encoding.
	UserNonVisibleData string `json:"userNonVisibleData,omitempty"`

	// Optional: If present, and set to “simpleMarkdownV1”, this parameter indicates that userVisibleData holds formatting characters which,
	// will potentially make the text displayed to the user nicer to look at.
	// For instructions check out https://www.bankid.com/utvecklare/guider/formatera-text
	UserVisibleDataFormat string `json:"userVisibleDataFormat,omitempty"`
}

func (r PhoneAuthRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type PhoneSignRequest struct {
	// Required: The user IP address as seen by RP. String. IPv4 and IPv6 is allowed.
	// Correct IP address must be the IP address representing the user agent (the end user device) as seen by the RP. In case of inbound proxy, special considerations may need to be taken into account to get the correct address.
	// In some use cases the IP address is not available, for instance in voice-based services. In these cases, the internal representation of those systems’ IP address may be used.
	EndUserIP string `json:"endUserIp"`

	// Required: Indicate if the user or the RP initiated the phone call.
	// 	- user: user called the RP
	// 	- RP: RP called the user
	CallInitiator string `json:"callInitiator"`

	// Optional: Requirements on how the auth order must be performed. See section Requirements below for more details.
	Requirement Requirement `json:"requirement,omitempty"`

	// Required: Text to be displayed to the user. String. The text can be formatted using CR, LF and CRLF for new lines.
	// The text must be encoded as UTF-8 and then base 64 encoded. 1 – 40,000 characters after base 64 encoding.
	UserVisibleData string `json:"userVisibleData,omitempty"`

	// Optional: Data is not displayed to the user. String. The value must be base 64-encoded. 1-1 500 characters after base 64-encoding.
	UserNonVisibleData string `json:"userNonVisibleData,omitempty"`

	// Optional: If present, and set to “simpleMarkdownV1”, this parameter indicates that userVisibleData holds formatting characters which,
	// will potentially make the text displayed to the user nicer to look at.
	// For instructions check out https://www.bankid.com/utvecklare/guider/formatera-text
	UserVisibleDataFormat string `json:"userVisibleDataFormat,omitempty"`
}

func (r PhoneSignRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Collects the result of a sign or auth order using orderRef as reference. RP should keep on calling collect every two seconds if status is pending. RP must abort if status indicates failed. The user identity is returned when complete.
type CollectRequest struct {
	// The orderRef returned from auth or sign.
	OrderRef string `json:"orderRef"`
}

func (r CollectRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Cancels an ongoing sign or auth order. This is typically used if the user cancels the order in your service or app.
type CancelRequest struct {
	// The orderRef returned from auth or sign.
	OrderRef string `json:"orderRef"`
}

func (r CancelRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// RP may use the requirement parameter to describe how a signature must be created and verified.
// A typical use case is to require Mobile BankID or a certain card reader.
// The following table describes requirements, their possible values and defaults.
type Requirement struct {
	// Users are required to sign the transaction with their PIN code, even if they have biometrics activated.
	// 	- Default: False, the user is not required to use pin code.
	Pincode bool `json:"pincode,omitempty"`

	// If present, and set to "true", the client needs to provide MRTD (Machine readable travel document) information to complete the order.
	// Only Swedish passports and national ID cards are supported.
	// 	- Default: The client does not need to provide MRTD information to complete the order.
	MRTD bool `json:"mrtd,omitempty"`

	// 	- "class1" (default) – The transaction must be performed using a card reader where the PIN code is entered on a computer keyboard, or a card reader of higher class.
	// 	- "class2" – The transaction must be performed using a card reader where the PIN code is entered on the reader, or a reader of higher class.
	// 	- "<"no value">" – defaults to "class1". This condition should be combined with a certificatePolicies for a smart card to avoid undefined behaviour.
	//  - Default: No card reader required.
	CardReader string `json:"cardReader,omitempty"`

	// The oid in certificate policies in the user certificate. List of String.
	// One wildcard ”” is allowed from position 5 and forward ie. 1.2.752.78.
	// The values for production BankIDs are:
	// 	- "1.2.752.78.1.1" - BankID on file
	// 	- "1.2.752.78.1.2" - BankID on smart card
	// 	- "1.2.752.78.1.5" - Mobile BankID
	// The values for test BankIDs are:
	// 	- "1.2.3.4.5" - BankID on file
	// 	- "1.2.3.4.10" - BankID on smart card
	// 	- "1.2.3.4.25" - Mobile BankID
	// 	- “1.2.752.60.1.6” - Test BankID for some BankID Banks
	// Default: If no set certificate policies, the following are default in the:
	// Production system
	// 	- 1.2.752.78.1.1
	// 	- 1.2.752.78.1.2
	// 	- 1.2.752.78.1.5
	// 	- 1.2.752.71.1.3
	// Test system:
	// 	- 1.2.3.4.5
	// 	- 1.2.3.4.10
	// 	- 1.2.3.4.25
	// 	- 1.2.752.60.1.6
	// 	- 1.2.752.71.1.3
	// If any certificate policy is set all default policies are dismissed.
	CertificatePolicies []string `json:"certificatePolicies,omitempty"`

	// A personal identification number to be used to complete the transaction.
	// If a BankID with another personal number attempts to sign the transaction, it fails.
	PersonalNumber string `json:"personalNumber,omitempty"`
	
	// Set the acceptable risk level for the transaction. If the risk is higher than the specified level, the transaction will be blocked.
	// String:
	// 	- low: only accept low risk orders
	// 	- moderate: accept low and moderate risk orders
	// The risk indication requires that the endUserIp is correct. An incorrect IP-address will result in legitimate transactions being blocked.
	// Default no risk level is set and the BankID service will not block any order regardless of the risk indication.
	Risk string `json:"risk,omitempty"`
}
