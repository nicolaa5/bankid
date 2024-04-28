package test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/nicolaa5/bankid/pkg/api"
	"github.com/nicolaa5/bankid/pkg/parameters"
	"github.com/nicolaa5/bankid/pkg/request"
	"github.com/stretchr/testify/require"
)

func TestAPI(t *testing.T) {
	cert, err := parameters.NewCert(
		parameters.WithPassphrase("qwerty"),
		parameters.WithSSLCertificatePath("../../certs/ssl_test.p12"),
		parameters.WithCACertificatePath("../../certs/ca_test.crt"),
	)
	require.NoError(t, err)

	b, err := api.New(parameters.Parameters{
		URL:  parameters.BankIDTestUrl,
		Cert: *cert,
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
				EndUserIP: randomIPv4(),
			})
			require.NoError(t, err)

			fmt.Println(response.OrderRef)
		})
	}
}

func randomIPv4() string {
	num := func() int { return 2 + rand.Intn(254) }
	return fmt.Sprintf("%d.%d.%d.%d", num(), num(), num(), num())
}
