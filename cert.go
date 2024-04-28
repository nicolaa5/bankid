package bankid

import (
	"fmt"
	"os"
)

// This certificate is used to authenticate the client to the BankID API.
type Cert struct {
	// Required: The password for your SSLCertificate
	Passphrase string `json:"passphrase"`

	// Required: Your organization's certificate signed by a trusted certificate authority (cert has .p12 extension).
	// Provided by the bank (the trusted CA) that you sign an agreement with, see https://www.bankid.com/en/foretag/kontakt-foeretag
	SSLCertificate []byte `json:"sslCertificate"`

	// Required: The BankID root certificate
	CACertificate []byte `json:"caCertificate"`
}

type CertPaths struct {
	// Required: The password for your SSLCertificate
	Passphrase string `json:"passphrase"`

	// Required: The path to your organization's certificate signed by a trusted certificate authority (cert has .p12 extension).
	SSLCertificatePath string `json:"sslCertificatePath"`

	// Required: The path to the BankID root certificate
	CACertificatePath string `json:"caCertificatePath"`
}

func CertFromPaths(params CertPaths) (*Cert, error) {
	cert := &Cert{}
	cert.Passphrase = params.Passphrase

	p12, err := os.ReadFile(params.SSLCertificatePath)
	if err != nil {
		return nil, fmt.Errorf("error reading .p12 file: %w", err)
	}

	cert.SSLCertificate = p12

	ca, err := os.ReadFile(params.CACertificatePath)
	if err != nil {
		return nil, fmt.Errorf("error reading root certificate file: %w", err)
	}

	cert.CACertificate = ca
	return cert, nil
}

func NewCert(opts ...CertOption) (*Cert, error) {
	cert := &Cert{}

	for _, opt := range opts {
		err := opt(cert)
		if err != nil {
			return nil, fmt.Errorf("invalid input: %w", err)
		}
	}

	return cert, nil
}

type CertOption func(*Cert) error

func WithSSLCertificate(cert []byte) CertOption {
	return func(c *Cert) error {
		c.SSLCertificate = cert
		return nil
	}
}

func WithPassphrase(passphrase string) CertOption {
	return func(c *Cert) error {
		c.Passphrase = passphrase
		return nil
	}
}

func WithCACertificate(cert []byte) CertOption {
	return func(c *Cert) error {
		c.CACertificate = cert
		return nil
	}
}

func WithSSLCertificatePath(path string) CertOption {
	return func(c *Cert) error {
		p12, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error reading .p12 file: %w", err)
		}

		c.SSLCertificate = p12
		return nil
	}
}

func WithCACertificatePath(path string) CertOption {
	return func(c *Cert) error {
		ca, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error reading root certificate file: %w", err)
		}

		c.CACertificate = ca
		return nil
	}
}
