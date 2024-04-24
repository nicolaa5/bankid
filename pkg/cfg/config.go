package cfg 

type Config struct {
	// The BankID API endpoint
	URL string

	// Your organization's certificate signed by a trusted certificate authority (cert has .p12 extension).
	// Provided by the bank (the trusted CA) that you sign an agreement with, see https://www.bankid.com/en/foretag/kontakt-foeretag
	SSLCertificate []byte

	// The password for your SSLCertificate
	Passphrase string

	// The BankID root certificate
	CACertificate []byte
}

type Option func(*Config)

func WithUrl(url string) Option {
	return func(c *Config) {
		c.URL = url
	}
}

func WithSSLCertificate(cert []byte) Option {
	return func(c *Config) {
		c.SSLCertificate = cert
	}
}

func WithPassphrase(passphrase string) Option {
	return func(c *Config) {
		c.Passphrase = passphrase
	}
}

func WithCACertificate(cert []byte) Option {
	return func(c *Config) {
		c.CACertificate = cert
	}
}

func NewConfig(opts ...Option) (*Config, error) {
	config := &Config{
		URL: "https://appapi2.bankid.com/rp/v6",
	}

	for _, opt := range opts {
		opt(config)
	}

	return config, nil 
}