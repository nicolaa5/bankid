package test

import (
	"fmt"
	"testing"

	"github.com/nicolaa5/bankid/pkg/api"
	"github.com/nicolaa5/bankid/pkg/parameters"
	"github.com/nicolaa5/bankid/pkg/request"
	"github.com/stretchr/testify/require"
)

func TestAPI(t *testing.T) {

	cert, err := parameters.NewCert(
		parameters.WithSSLCertificatePath("certs/client.p12"),
		parameters.WithCACertificatePath("certs/ca.pem"),
	)
	require.NoError(t, err)

	b, err := api.New(parameters.Parameters{
		URL:  parameters.BankIDTestUrl,
		Cert: cert,
	})
	require.NoError(t, err)
	
	for _, tt := range []struct {
		name string
	}{
		{
			name: "Authenticate",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			response, err := b.Auth(request.AuthRequest{
				EndUserIP: "",
			})
			require.NoError(t, err)

			fmt.Println(response.OrderRef)
		})
	}
}
