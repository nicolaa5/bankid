package validate

import (
	"fmt"

	"github.com/nicolaa5/bankid/pkg/parameters"
)

func Config(params *parameters.Parameters) error {
	if params.SSLCertificate == nil && params.SSLCertificatePath == "" {
		return fmt.Errorf("ssl certificate is not provided")
	}

	if params.CACertificate == nil && params.CACertificatePath == "" {
		return fmt.Errorf("ca root certificate is not provided")
	}

	// Set the URL to the default production endpoint if not provided
	if params.URL == "" {
		params.URL = "https://appapi2.bankid.com/rp/v6"
	}

	// Set the timeout to 5 seconds if not provided
	if params.Timeout == 0 {
		params.Timeout = 5
	}

	return nil
}
