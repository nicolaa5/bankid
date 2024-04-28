package test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/nicolaa5/bankid"
	"github.com/stretchr/testify/require"
)

func TestCertPaths(t *testing.T) {
	_, err := bankid.CertFromPaths(bankid.CertPaths{
		Passphrase:         bankid.BankIDTestPassphrase,
		SSLCertificatePath: "../../certs/ssl_test.p12",
		CACertificatePath:  "../../certs/ca_test.crt",
	})
	require.NoError(t, err)
}

func TestAuthenticate(t *testing.T) {
	b, err := bankid.New(bankid.Parameters{
		URL: bankid.BankIDTestUrl,
		Cert: bankid.Cert{
			Passphrase:     bankid.BankIDTestPassphrase,
			SSLCertificate: bankid.SSLTestCertificate,
			CACertificate:  bankid.CATestCertificate,
		},
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
			})

			require.NoError(t, err)
			require.NotEmpty(t, response.OrderRef)
			require.NotEmpty(t, response.AutoStartToken)
			require.NotEmpty(t, response.QrStartSecret)
			require.NotEmpty(t, response.QrStartToken)
		})
	}
}

func randomIPv4() string {
	num := func() int { return 2 + rand.Intn(254) }
	return fmt.Sprintf("%d.%d.%d.%d", num(), num(), num(), num())
}
