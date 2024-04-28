package parameters

import (
	"fmt"
	"os"
)

// This certificate is used to authenticate the client to the BankID API.
type Cert struct {
	// Required: The password for your SSLCertificate
	Passphrase string `json:"passphrase"`

	// Your organization's certificate signed by a trusted certificate authority (cert has .p12 extension).
	// Provided by the bank (the trusted CA) that you sign an agreement with, see https://www.bankid.com/en/foretag/kontakt-foeretag
	SSLCertificate []byte `json:"sslCertificate"`

	// The BankID root certificate
	CACertificate []byte `json:"caCertificate"`
}

func NewCert(opts ...CertOption) (Cert, error) {
	config := Cert{}

	for _, opt := range opts {
		opt(config)
	}

	return config, nil
}

type CertOption func(Cert) error

func WithSSLCertificate(cert []byte) CertOption {
	return func(c Cert) error {
		c.SSLCertificate = cert
		return nil
	}
}

func WithPassphrase(passphrase string) CertOption {
	return func(c Cert) error {
		c.Passphrase = passphrase
		return nil
	}
}

func WithCACertificate(cert []byte) CertOption {
	return func(c Cert) error {
		c.CACertificate = cert
		return nil
	}
}

func WithSSLCertificatePath(path string) CertOption {
	return func(c Cert) error {
		p12, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error reading .p12 file: %w", err)
		}
		
		c.SSLCertificate = p12
		return nil
	}
}

func WithCACertificatePath(path string) CertOption {
	return func(c Cert) error {
		ca, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error reading root certificate file: %w", err)
		}

		c.CACertificate = ca
		return nil
	}
}
