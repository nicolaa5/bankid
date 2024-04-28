package test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/nicolaa5/bankid"
	"github.com/stretchr/testify/require"
)

func TestAPI(t *testing.T) {
	cert, err := bankid.CertFromPaths(bankid.CertPaths{
		Passphrase:         bankid.BankIDTestPassphrase,
		SSLCertificatePath: "../../certs/ssl_test.p12",
		CACertificatePath:  "../../certs/ca_test.crt",
	})
	require.NoError(t, err)

	b, err := bankid.New(bankid.Parameters{
		URL:  bankid.BankIDTestUrl,
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

			response, err := b.Auth(bankid.AuthRequest{
				EndUserIP: randomIPv4(),
				Requirement: bankid.Requirement{
					PersonalNumber: "199207201337",
				},
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
