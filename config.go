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
func (c *Config) EnsureRequired() {
	if c.URL == "" {
		// Set the URL to the default production endpoint if not provided
		c.URL = BankIDURL
	}

	// Set the timeout to 5 seconds if not provided
	if c.Timeout == 0 {
		c.Timeout = 5
	}

	// Use the BankID CA root certificate by default for requests 
	if c.CACertificate == nil && c.URL == BankIDURL {
		c.CACertificate = CAProdCertificate
	}
}