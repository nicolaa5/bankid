package bankid

import (
	"fmt"
	"net"
	"unicode/utf8"

	personnummer "github.com/personnummer/go/v3"
)

type ValidateOption func() error

// Ensure that the requirements that are required are present and in the expected format
func validate(opts ...ValidateOption) error {
	for _, opt := range opts {
		err := opt()
		if err != nil {
			return err
		}
	}

	return nil
}

func validateRequired(body RequestBody) ValidateOption {
	return func() error {
		switch v := body.(type) {
		case AuthRequest:
			if v.EndUserIP == "" {
				return RequiredInputMissingError{Message: fmt.Sprintf("EndUserIP is missing but required by BankID for request: %T", v)}
			}
		case SignRequest:
			if v.EndUserIP == "" {
				return RequiredInputMissingError{Message: fmt.Sprintf("EndUserIP is missing but required by BankID for request: %T", v)}
			}

			if v.UserVisibleData == "" || !utf8.ValidString(v.UserVisibleData) {
				return RequiredInputMissingError{Message: fmt.Sprintf("UserVisibleData is missing or invalid, it's required by BankID for request: %T", v)}
			}

		case PhoneAuthRequest:
			if v.PersonalNumber == "" {
				return RequiredInputMissingError{Message: fmt.Sprintf("PersonalNumber is missing but required by BankID for request: %T", v)}
			}

			if v.CallInitiator == "" {
				return RequiredInputMissingError{Message: fmt.Sprintf("CallInitiator is missing but required by BankID for request: %T", v)}
			}

		case PhoneSignRequest:
			if v.PersonalNumber == "" {
				return RequiredInputMissingError{Message: fmt.Sprintf("PersonalNumber is missing but required by BankID for request: %T", v)}
			}

			if v.CallInitiator == "" {
				return RequiredInputMissingError{Message: fmt.Sprintf("CallInitiator is missing but required by BankID for request: %T", v)}
			}

			if v.UserVisibleData == "" || !utf8.ValidString(v.UserVisibleData) {
				return RequiredInputMissingError{Message: fmt.Sprintf("UserVisibleData is missing or invalid, it's required by BankID for request: %T", v)}
			}

		}
		return nil
	}
}

func validateEndUserIP(endUserIP string) ValidateOption {
	return func() error {
		if valid := isValidIP(endUserIP); !valid {
			return InputInvalidError{Message: fmt.Sprintf("EndUserIP: %s is invalid", endUserIP)}
		}

		return nil
	}
}

func validateCallInitiator(callInitiator string) ValidateOption {
	return func() error {
		if callInitiator != "user" && callInitiator != "RP" {
			return InputInvalidError{Message: fmt.Sprint("CallInitator is not formatted correctly, it should be 'user' or 'RP'")}
		}

		return nil
	}
}

func validateCardReader(cardReader string) ValidateOption {
	return func() error {
		if cardReader != "" && cardReader != "class1" && cardReader != "class2" {
			return InputInvalidError{Message: fmt.Sprintf("CardReader input is invalid, it should be 'class1' or 'class2'")}
		}

		return nil
	}
}

func validateCertificatePolicies(certificatePolicies []string) ValidateOption {
	validOIDs := []string{
		"1.2.752.78.1.1", //  Represents BankID on file
		"1.2.752.78.1.2", //  Represents BankID on smart card
		"1.2.752.78.1.5", //  Represents Mobile BankID
		"1.2.3.4.5",      //  Test BankID on file
		"1.2.3.4.10",     //  Test BankID on smart card
		"1.2.3.4.25",     //  Test Mobile BankID
		"1.2.752.60.1.6", //  Test BankID for certain BankID Banks
	}

	return func() error {
		for _, oid := range certificatePolicies {
			var valid bool
			for _, validOID := range validOIDs {
				if validOID == oid {
					valid = true
				}
			}

			if !valid {
				return InputInvalidError{Message: fmt.Sprintf("Certificate Policy input: %s is invalid, check BankID for valid certificate policies", oid)}
			}
		}

		return nil
	}
}

func validatePersonalNumber(personalNumber string) ValidateOption {
	return func() error {
		if personalNumber != "" && !personnummer.Valid(personalNumber) {
			return InputInvalidError{Message: fmt.Sprintf("Personnummer: %s is not formatted correctly", personalNumber)}
		}
		return nil
	}
}

func validateRequirement(requirement *Requirement) ValidateOption {
	return func() error {
		// bankid request requirements are optional so nil is accepted
		if requirement == nil {
			return nil
		}

		opts := []ValidateOption{
			validatePersonalNumber(requirement.PersonalNumber),
			validateCertificatePolicies(requirement.CertificatePolicies),
			validateCardReader(requirement.CardReader),
		}

		for _, opt := range opts {
			err := opt()
			if err != nil {
				return err
			}
		}

		return nil
	}
}

func validateParameters(p Config) ValidateOption {
	return func() error {

		if p.SSLCertificate == nil {
			return InputInvalidError{Message: fmt.Sprint("ssl certificate is not provided")}
		}

		if p.CACertificate == nil {
			return InputInvalidError{Message: fmt.Sprint("ca root certificate is not provided")}
		}

		return nil
	}
}

func isValidIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil
}
