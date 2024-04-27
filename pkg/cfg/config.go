package cfg

import "fmt"

type Config struct {
	Cert

	// The URL to BankID API
	URL string `json:"url"`

	// The timeout for the request to BankID API in seconds.
	Timeout int `json:"timeout"`
}

func (c Config) Validate() error {
	if c.URL == "" {
		// if url is not set we assume the default value
		c.URL = "https://appapi2.bankid.com/rp/v6"
	}

	if c.SSLCertificate == nil && c.SSLCertificatePath == "" {
		return fmt.Errorf("ssl certificate is not provided")
	}
	if c.CACertificate == nil && c.CACertificatePath == "" {
		return fmt.Errorf("ca root certificate is not provided")
	}
	return nil
}
