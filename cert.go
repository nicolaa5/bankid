package bankid

import (
	_ "embed"
	"fmt"
)

var (
	//go:embed certs/ca_prod.crt
	CAProdCertificate []byte

	//go:embed certs/ca_test.crt
	CATestCertificate []byte

	//go:embed certs/FPTestcert5_20240610.p12
	P12TestCertificate []byte

	//go:embed certs/FPTestcert5_20240610.pem
	PEMTestCertificate string
)

// The certificate is used to authenticate the RP service to the BankID API.
type Certificate interface {
	Validate() error
	CA() []byte
}

type P12Cert struct {
	// Required: Your organization's certificate signed by a trusted certificate authority (cert has .p12 extension).
	// Provided by the bank (the trusted CA) that you sign an agreement with, see https://www.bankid.com/en/foretag/kontakt-foeretag
	Certificate []byte `json:"certificate"`

	// Required: The password for your P12Certificate
	Passphrase string `json:"passphrase"`

	// Optional: A CA root certificate. This lib uses the BankID root certificate as the default
	CACertificate []byte `json:"caCertificate"`
}

func (c P12Cert) Validate() error {
	if c.Certificate == nil {
		return fmt.Errorf("p12 certificate is not provided")
	}

	if c.Passphrase == "" {
		return fmt.Errorf("passphrase for p12 certificate is not provided")
	}

	return nil
}

func (c P12Cert) CA() []byte {
	return c.CACertificate
}

type PEMCert struct {
	Certificate string `json:"certificate"`

	// Required: The password for your .pem Certificate
	Passphrase string `json:"passphrase"`

	// Optional: A CA root certificate. This lib uses the BankID root certificate as the default
	CACertificate []byte `json:"caCertificate"`
}

func (c PEMCert) Validate() error {
	if c.Certificate == "" {
		return fmt.Errorf("private key is not provided")
	}

	if c.Passphrase == "" {
		return fmt.Errorf("passphrase is not provided")
	}

	return nil
}

func (c PEMCert) CA() []byte {
	return c.CACertificate
}
