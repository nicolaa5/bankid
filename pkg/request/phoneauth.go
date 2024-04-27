package request

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
