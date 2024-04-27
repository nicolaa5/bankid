package request

// RP may use the requirement parameter to describe how a signature must be created and verified.
// A typical use case is to require Mobile BankID or a certain card reader.
// The following table describes requirements, their possible values and defaults.
type Requirement struct {
	// Users are required to sign the transaction with their PIN code, even if they have biometrics activated.
	// 	- Default: False, the user is not required to use pin code.
	Pincode bool `json:"pincode"`

	// If present, and set to "true", the client needs to provide MRTD (Machine readable travel document) information to complete the order.
	// Only Swedish passports and national ID cards are supported.
	// 	- Default: The client does not need to provide MRTD information to complete the order.
	MRTD bool `json:"mrtd"`

	// 	- "class1" (default) – The transaction must be performed using a card reader where the PIN code is entered on a computer keyboard, or a card reader of higher class.
	// 	- "class2" – The transaction must be performed using a card reader where the PIN code is entered on the reader, or a reader of higher class.
	// 	- "<"no value">" – defaults to "class1". This condition should be combined with a certificatePolicies for a smart card to avoid undefined behaviour.
	//  - Default: No card reader required.
	CardReader string `json:"cardReader"`

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
	CertificatePolicies []string `json:"certificatePolicies"`

	// A personal identification number to be used to complete the transaction.
	// If a BankID with another personal number attempts to sign the transaction, it fails.
	PersonalNumber string `json:"personalNumber"`
}
