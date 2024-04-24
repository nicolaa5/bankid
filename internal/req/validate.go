package req

import (
	"fmt"

	"github.com/nicolaa5/bankid/pkg/cfg"
)

func validateRequest(config cfg.Config) error {
    if config.URL == "" {
		config.URL = "https://appapi2.bankid.com/rp/v6"
        return fmt.Errorf("url is not provided")
    }
}