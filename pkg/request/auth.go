package request

import "encoding/json"

type AuthRequest struct {
	// Required: The user IP address as seen by RP. String. IPv4 and IPv6 is allowed.
	// Correct IP address must be the IP address representing the user agent (the end user device) as seen by the RP. In case of inbound proxy, special considerations may need to be taken into account to get the correct address.
	// In some use cases the IP address is not available, for instance in voice-based services. In these cases, the internal representation of those systems’ IP address may be used.
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
}

func (r AuthRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
