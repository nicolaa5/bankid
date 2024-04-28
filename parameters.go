package bankid

import "fmt"

const (
	BankIDURL            = "https://appapi2.bankid.com/rp/v6.0"
	BankIDTestUrl        = "https://appapi2.test.bankid.com/rp/v6.0"
	BankIDTestPassphrase = "qwerty123"
)

type Parameters struct {
	// Required: The SSL & CA certificate for the client.
	Certificate

	// Optional: The URL to BankID API, can be set to the test or production endpoint.
	// Default: "https://appapi2.bankid.com/rp/v6.0"
	URL string `json:"url"`

	// Optional: The timeout for the request to BankID API in seconds.
	// Default: 5
	Timeout int `json:"timeout"`
}

func (p Parameters) Validate() error {
	if p.URL == "" {
		// Set the URL to the default production endpoint if not provided
		p.URL = BankIDURL
	}

	if p.SSLCertificate == nil {
		return fmt.Errorf("ssl certificate is not provided")
	}

	if p.CACertificate == nil {
		return fmt.Errorf("ca root certificate is not provided")
	}

	// Set the timeout to 5 seconds if not provided
	if p.Timeout == 0 {
		p.Timeout = 5
	}

	return nil
}
