package test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/google/uuid"
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
			name: "Authenticate, collect, cancel",
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

func TestErrorCodes(t *testing.T) {
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
		expected bankid.BankIDError
		f        func()
	}{
		{expected: bankid.ErrInvalidParameters},
		{expected: bankid.ErrAlreadyInProgress},
		{expected: bankid.ErrUnauthorized},
	} {
		t.Run("Test Error Codes", func(t *testing.T) {
			test := tt
			t.Parallel()

			switch test.expected.ErrorCode {
			case bankid.InvalidParameters:
				//empty invalid IP as EndUserIP
				_, err := b.Auth(bankid.AuthRequest{
					EndUserIP: "",
				})
				require.Error(t, err)

				e, ok := err.(bankid.BankIDError)
				require.True(t, ok)
				require.Equal(t, e.ErrorCode, bankid.InvalidParameters)

				//using non-existent orderRef
				uniqueID := uuid.New()
				_, err = b.Collect(bankid.CollectRequest{
					OrderRef: uniqueID.String(),
				})
				require.Error(t, err)

				e, ok = err.(bankid.BankIDError)
				require.True(t, ok)
				require.Equal(t, e.ErrorCode, bankid.InvalidParameters)

			case bankid.AlreadyInProgress:
				ip := randomIPv4()
				personNummer := randomPersonnummer()

				//first request
				_, err := b.Auth(bankid.AuthRequest{
					EndUserIP: ip,
					Requirement: bankid.Requirement{
						PersonalNumber: personNummer,
					},
				})
				require.NoError(t, err)

				//second request with the same personNummer
				_, err = b.Auth(bankid.AuthRequest{
					EndUserIP: ip,
					Requirement: bankid.Requirement{
						PersonalNumber: personNummer,
					},
				})
				require.Error(t, err)
				fmt.Printf("%v", err.Error())

				e, ok := err.(bankid.BankIDError)
				require.True(t, ok)
				require.Equal(t, e.ErrorCode, bankid.AlreadyInProgress)

			case bankid.NotFound:
				b, err := bankid.New(bankid.Parameters{
					// add non-existent path to URL
					URL: "https://appapi2.test.bankid.com/rp/v6.0/forbidden/path",
					Certificate: bankid.Certificate{
						Passphrase:     bankid.BankIDTestPassphrase,
						SSLCertificate: bankid.SSLTestCertificate,
						CACertificate:  bankid.CATestCertificate,
					},
				})
				require.NoError(t, err)

				_, err = b.Auth(bankid.AuthRequest{
					EndUserIP: randomIPv4(),
				})
				require.Error(t, err)

				e, ok := err.(bankid.BankIDError)
				require.True(t, ok)
				require.Equal(t, e.ErrorCode, bankid.ErrUnauthorized)
			}
		})
	}
}

func randomIPv4() string {
	num := func() int { return 2 + rand.Intn(254) }
	return fmt.Sprintf("%d.%d.%d.%d", num(), num(), num(), num())
}

func randomNumber(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func randomPersonnummer() string {
	year := randomNumber(1900, 2024)
	month := randomNumber(1, 12)
	day := randomNumber(1, 28)
	serialNumber := randomNumber(1000, 9999)

	return fmt.Sprintf("%04d%02d%02d%04d", year, month, day, serialNumber)
}
