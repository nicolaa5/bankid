package parameters

type Cert struct {
	// Required: The password for your SSLCertificate
	Passphrase string `json:"passphrase"`

	// Your organization's certificate signed by a trusted certificate authority (cert has .p12 extension).
	// Provided by the bank (the trusted CA) that you sign an agreement with, see https://www.bankid.com/en/foretag/kontakt-foeretag
	SSLCertificate []byte `json:"sslCertificate"`

	// The path to your SSLCertificate
	SSLCertificatePath string `json:"sslCertificatePath"`

	// The BankID root certificate
	CACertificate []byte `json:"caCertificate"`

	// The path to your CACertificate
	CACertificatePath string `json:"caCertificatePath"`
}

func NewCert(opts ...CertOption) (Cert, error) {
	config := Cert{}

	for _, opt := range opts {
		opt(config)
	}

	return config, nil
}

func NewCertFromPaths(opts ...CertPathOption) Cert {
	config := Cert{}

	for _, opt := range opts {
		opt(config)
	}

	return config
}

type CertOption func(Cert)
type CertPathOption func(Cert)

func WithSSLCertificate(cert []byte) CertOption {
	return func(c Cert) {
		c.SSLCertificate = cert
	}
}

func WithPassphrase(passphrase string) CertOption {
	return func(c Cert) {
		c.Passphrase = passphrase
	}
}

func WithCACertificate(cert []byte) CertOption {
	return func(c Cert) {
		c.CACertificate = cert
	}
}

func WithSSLCertificatePath(path string) CertPathOption {
	return func(c Cert) {
		c.SSLCertificatePath = path
	}
}

func WithCACertificatePath(path string) CertPathOption {
	return func(c Cert) {
		c.CACertificatePath = path
	}
}
