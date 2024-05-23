package cli

import (
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
)

func endUserIpPrompt() *survey.Input {
	return &survey.Input{
		Message: "Please provide the IP address of the Enduser",
		Help:    "EndUserIP is required by BankID because it allows you to improve security by allowing you to compare the enduser's IP you receive with the IP send in the completiondata by BankID.",
		Default: "192.168.0.0",
		Suggest: func(toComplete string) []string {
			return []string{"192.168.0.0", "127.0.0.0"}
		},
	}
}

func passphrasePrompt() *survey.Password {
	return &survey.Password{
		Message: "Please provide the Passphrase with which your SSLCertificate is encrypted",
		Help:    "The password that is required is used to  protect and secure the private key with which the CSR is signed that BankID uses to generate your SSLCertificate",
	}
}

func sslCertPrompt() *survey.Input {
	return &survey.Input{
		Message: "Please provide the path to your SSLCertificate",
		Help:    "SSL Certificates help establish trust between your company and BankID by providing a secure connection and validating your organization's identity",
		Suggest: func(toComplete string) []string {
			files, _ := filepath.Glob(toComplete + "*")
			return files
		},
	}
}

func userDataVisibleDataPrompt() *survey.Input {
	return &survey.Input{
		Message: "Provide the text that's visible to the user when signing",
		Help:    "The visible data helps the enduser understand what they are signing",
	}
}