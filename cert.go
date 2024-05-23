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

	//go:embed certs/FPTestcert4_20230629.p12
	SSLTestCertificate []byte
)

// This certificate is used to authenticate the RP service to the BankID API.
type Certificate struct {
	// Required: The password for your SSLCertificate
	Passphrase string `json:"passphrase"`

	// Required if SSLCertificatePath is not provided: Your organization's certificate signed by a trusted certificate authority (cert has .p12 extension).
	// Provided by the bank (the trusted CA) that you sign an agreement with, see https://www.bankid.com/en/foretag/kontakt-foeretag
	SSLCertificate []byte `json:"sslCertificate"`

	// Optional: A CA root certificate. This lib has uses the BankID root certificate as the default
	CACertificate []byte `json:"caCertificate"`

	// Required if SSLCertificate is not provided: The path to your organization's certificate signed by a trusted certificate authority (cert has .p12 extension).
	SSLCertificatePath string `json:"sslCertificatePath"`

	// Optional: The path to the BankID root certificate
	CACertificatePath string `json:"caCertificatePath"`
}

func (c Certificate) Validate() error {
	if c.SSLCertificate == nil && c.SSLCertificatePath == "" {
		return fmt.Errorf("ssl certificate and path are not provided")
	}

	if c.Passphrase == ""{
		return fmt.Errorf("passphrase for ssl certificate is not provided")
	}

	return nil
}

func CertificateFromPath(path string) ([]byte, error) {
	cert, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading .p12 file: %w", err)
	}

	return cert, nil
}
