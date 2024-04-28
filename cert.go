package bankid

import (
	_ "embed"
	"fmt"
	"os"
)

var (
	//go:embed certs/ca_prod.crt
	CAProdCertificate []byte

	//go:embed certs/ca_test.crt
	CATestCertificate []byte

	//go:embed certs/ssl_test.p12
	SSLTestCertificate []byte
)

// This certificate is used to authenticate the client to the BankID API.
type Certificate struct {
	// Required: The password for your SSLCertificate
	Passphrase string `json:"passphrase"`

	// Required: Your organization's certificate signed by a trusted certificate authority (cert has .p12 extension).
	// Provided by the bank (the trusted CA) that you sign an agreement with, see https://www.bankid.com/en/foretag/kontakt-foeretag
	SSLCertificate []byte `json:"sslCertificate"`

	// Required: The BankID root certificate
	CACertificate []byte `json:"caCertificate"`
}

type CertificatePaths struct {
	// Required: The password for your SSLCertificate
	Passphrase string `json:"passphrase"`

	// Required: The path to your organization's certificate signed by a trusted certificate authority (cert has .p12 extension).
	SSLCertificatePath string `json:"sslCertificatePath"`

	// Required: The path to the BankID root certificate
	CACertificatePath string `json:"caCertificatePath"`
}

func CertificateFromPaths(params CertificatePaths) (*Certificate, error) {
	if params.Passphrase == "" {
		return nil, fmt.Errorf("passphrase is required")
	}

	if params.SSLCertificatePath == "" {
		return nil, fmt.Errorf("ssl certificate path is required")
	}

	if params.CACertificatePath == "" {
		return nil, fmt.Errorf("ca certificate path is required")
	}

	cert := &Certificate{}
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
