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
)

// The certificate is used to authenticate the RP service to the BankID API.
type Certificate interface {
	String() string
	Validate() error
	CA() []byte
}

type P12Cert struct {
	// Required: The password for your SSLCertificate
	Passphrase string `json:"passphrase"`

	// Required: Your organization's certificate signed by a trusted certificate authority (cert has .p12 extension).
	// Provided by the bank (the trusted CA) that you sign an agreement with, see https://www.bankid.com/en/foretag/kontakt-foeretag
	P12Certificate []byte `json:"p12Certificate"`

	// Optional: A CA root certificate. This lib uses the BankID root certificate as the default
	CACertificate []byte `json:"caCertificate"`
}

func (c P12Cert) String() string {
	return string(c.P12Certificate)
}

func (c P12Cert) Validate() error {
	if c.P12Certificate == nil {
		return fmt.Errorf("ssl certificate is not provided")
	}

	if c.Passphrase == "" {
		return fmt.Errorf("passphrase for ssl certificate is not provided")
	}

	return nil
}

func (c P12Cert) CA() []byte {
	return c.CACertificate
}

type PEMCert struct {
	PublicKey []byte `json:"publicKey"`

	PrivateKey []byte `json:"privateKey"`

	// Optional: A CA root certificate. This lib uses the BankID root certificate as the default
	CACertificate []byte `json:"caCertificate"`
}

func (c PEMCert) String() string {
	return string(c.PublicKey)
}

func (c PEMCert) Validate() error {
	if c.PrivateKey == nil {
		return fmt.Errorf("private key is not provided")
	}

	if c.PublicKey == nil {
		return fmt.Errorf("public key is not provided")
	}

	return nil
}

func (c PEMCert) CA() []byte {
	return c.CACertificate
}