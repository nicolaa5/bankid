package parameters

import "fmt"

const (
	BankIDURL     = "https://appapi2.bankid.com/rp/v6.0"
	BankIDTestUrl = "https://appapi2.test.bankid.com/rp/v6.0"
)

type Parameters struct {
	Cert

	// Optional: The URL to BankID API
	// Default: "https://appapi2.bankid.com/rp/v6.0"
	URL string `json:"url"`

	// Optional: The timeout for the request to BankID API in seconds.
	// Default: 5
	Timeout int `json:"timeout"`
}

func (c Parameters) Validate() error {
	if c.URL == "" {
		// if url is not set we assume the default value
		c.URL = BankIDURL
	}

	if c.SSLCertificate == nil && c.SSLCertificatePath == "" {
		return fmt.Errorf("ssl certificate is not provided")
	}
	if c.CACertificate == nil && c.CACertificatePath == "" {
		return fmt.Errorf("ca root certificate is not provided")
	}
	return nil
}
