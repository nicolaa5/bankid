package cfg

type Cert struct {
	// Your organization's certificate signed by a trusted certificate authority (cert has .p12 extension).
	// Provided by the bank (the trusted CA) that you sign an agreement with, see https://www.bankid.com/en/foretag/kontakt-foeretag
	SSLCertificate []byte `json:"sslCertificate"`

	// The path to your SSLCertificate
	SSLCertificatePath string `json:"sslCertificatePath"`

	// The password for your SSLCertificate
	Passphrase string `json:"passphrase"`

	// The BankID root certificate
	CACertificate []byte `json:"caCertificate"`

	// The path to your CACertificate
	CACertificatePath string `json:"caCertificatePath"`
}

type Option func(*Cert)

func WithSSLCertificate(cert []byte) Option {
	return func(c *Cert) {
		c.SSLCertificate = cert
	}
}

func WithPassphrase(passphrase string) Option {
	return func(c *Cert) {
		c.Passphrase = passphrase
	}
}

func WithCACertificate(cert []byte) Option {
	return func(c *Cert) {
		c.CACertificate = cert
	}
}

func NewConfig(opts ...Option) (*Cert, error) {
	config := &Cert{}

	for _, opt := range opts {
		opt(config)
	}

	return config, nil
}
