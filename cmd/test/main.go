package test

import (
	"testing"

	"github.com/nicolaa5/bankid/pkg/api"
	"github.com/nicolaa5/bankid/pkg/parameters"
	"github.com/stretchr/testify/require"
)

func TestAPI(t *testing.T) {
	for _, tt := range []struct {
		name string
	}{
		{
			name: "Authenticate",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cert, err := parameters.NewCert(
				parameters.WithSSLCertificate([]byte("certs/client.p12")),
				parameters.WithCACertificate([]byte("certs/ca.pem")),
			)
			require.NoError(t, err)

			config := parameters.Parameters{
				URL: parameters.BankIDTestUrl,
				Cert: *cert,
			}

			api.New(config)
		})
	}
}
