package bankid

type Config struct {
	// Required: The SSL & CA certificate for the client.
	Certificate

	// Optional: The URL to BankID API, can be set to the test or production endpoint.
	// Default: "https://appapi2.bankid.com/rp/v6.0"
	URL string `json:"url"`

	// Optional: The timeout for the request to BankID API in seconds.
	// Default: 5
	Timeout int `json:"timeout"`
}

// Ensures input data is set based on BankID requirements or leaves the input unchanged if it's valid or optional
func (c *Config) UseDefault() {
	if c.URL == "" {
		// Set the URL to the default production endpoint if not provided
		c.URL = BankIDURL
	}

	// Set the timeout to 5 seconds if not provided
	if c.Timeout == 0 {
		c.Timeout = 5
	}

	// Use the BankID CA root certificate for dev and prod scnearios by default for requests
	if c.CA() == nil {
		switch v := c.Certificate.(type) {
		case P12Cert:
			if c.URL == BankIDURL {
				v.CACertificate = CAProdCertificate
			} else if c.URL == BankIDTestUrl {
				v.CACertificate = CATestCertificate
			}
		case PEMCert:
			if c.URL == BankIDURL {
				v.CACertificate = CAProdCertificate
			} else if c.URL == BankIDTestUrl {
				v.CACertificate = CATestCertificate
			}
		}
	}
}
