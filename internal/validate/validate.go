package validate

import (
	"fmt"

	"github.com/nicolaa5/bankid/pkg/cfg"
)

func Config(config *cfg.Config) error {
	if config.SSLCertificate == nil && config.SSLCertificatePath == "" {
		return fmt.Errorf("ssl certificate is not provided")
	}

	if config.URL == "" {
		config.URL = "https://appapi2.bankid.com/rp/v6"
		return fmt.Errorf("url is not provided")
	}

	return nil
}
