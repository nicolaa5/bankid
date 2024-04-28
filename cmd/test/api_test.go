package test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/nicolaa5/bankid"
	"github.com/stretchr/testify/require"
)

func TestCertPaths(t *testing.T) {
	_, err := bankid.CertificateFromPaths(bankid.CertificatePaths{
		Passphrase:         bankid.BankIDTestPassphrase,
		SSLCertificatePath: "../../certs/ssl_test.p12",
		CACertificatePath:  "../../certs/ca_test.crt",
	})
	require.NoError(t, err)
}

func TestAuthenticate(t *testing.T) {
	b, err := bankid.New(bankid.Parameters{
		URL: bankid.BankIDTestUrl,
		Certificate: bankid.Certificate{
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

			authResponse, err := b.Auth(bankid.AuthRequest{
				EndUserIP: randomIPv4(),
			})

			require.NoError(t, err)
			require.NotEmpty(t, authResponse.OrderRef)
			require.NotEmpty(t, authResponse.AutoStartToken)
			require.NotEmpty(t, authResponse.QrStartSecret)
			require.NotEmpty(t, authResponse.QrStartToken)

			collectResponse, err := b.Collect(bankid.CollectRequest{
				OrderRef: authResponse.OrderRef,
			})

			require.NoError(t, err)
			require.Equal(t, collectResponse.OrderRef, authResponse.OrderRef)
			require.Equal(t, collectResponse.Status, bankid.Pending)
			require.Equal(t, collectResponse.HintCode, bankid.OutstandingTransaction)

			_, err = b.Cancel(bankid.CancelRequest{
				OrderRef: authResponse.OrderRef,
			})

			require.NoError(t, err)
		})
	}
}

func randomIPv4() string {
	num := func() int { return 2 + rand.Intn(254) }
	return fmt.Sprintf("%d.%d.%d.%d", num(), num(), num(), num())
}
