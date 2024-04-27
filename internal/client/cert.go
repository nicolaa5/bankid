package client

import (
	"fmt"
	"os"

	"github.com/nicolaa5/bankid/pkg/parameters"
)

func readCert(config *parameters.Parameters) error {
	if config.SSLCertificate == nil && config.SSLCertificatePath != "" {
		p12, err := os.ReadFile(config.SSLCertificatePath)
		if err != nil {
			return fmt.Errorf("error reading .p12 file: %w", err)
		}

		config.SSLCertificate = p12
	}

	if config.CACertificate == nil && config.CACertificatePath != "" {
		ca, err := os.ReadFile(config.CACertificatePath)
		if err != nil {
			return fmt.Errorf("error reading root certificate file: %w", err)
		}

		config.CACertificate = ca
	}

	return nil
}
