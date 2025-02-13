package test

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/nicolaa5/bankid"
	"github.com/stretchr/testify/require"
)

func TestAuthenticate(t *testing.T) {
	ctx := context.Background()

	p12Config, err := bankid.New(bankid.Config{
		URL: bankid.BankIDTestUrl,
		Certificate: bankid.P12Cert{
			Passphrase:    bankid.BankIDTestPassphrase,
			Certificate:   bankid.P12TestCertificate,
			CACertificate: bankid.CATestCertificate,
		},
	})
	require.NoError(t, err)

	pemConfig, err := bankid.New(bankid.Config{
		URL: bankid.BankIDTestUrl,
		Certificate: bankid.PEMCert{
			Certificate:   bankid.PEMTestCertificate,
			Passphrase:    bankid.BankIDTestPassphrase,
			CACertificate: bankid.CATestCertificate,
		},
	})
	require.NoError(t, err)

	for _, tt := range []struct {
		name   string
		client bankid.BankID
	}{
		{
			name:   "Authenticate, collect, cancel - using .p12 cert",
			client: p12Config,
		},
		{
			name:   "Authenticate, collect, cancel - using .pem cert",
			client: pemConfig,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			authResponse, err := tt.client.Auth(ctx, bankid.AuthRequest{
				EndUserIP: randomIPv4(),
			})

			require.NoError(t, err)
			require.NotEmpty(t, authResponse.OrderRef)
			require.NotEmpty(t, authResponse.AutoStartToken)
			require.NotEmpty(t, authResponse.QrStartSecret)
			require.NotEmpty(t, authResponse.QrStartToken)

			collectResponse, err := tt.client.Collect(ctx, bankid.CollectRequest{
				OrderRef: authResponse.OrderRef,
			})

			require.NoError(t, err)
			require.Equal(t, collectResponse.OrderRef, authResponse.OrderRef)
			require.Equal(t, collectResponse.Status, bankid.Pending)
			require.Equal(t, collectResponse.HintCode, bankid.OutstandingTransaction)

			_, err = tt.client.Cancel(ctx, bankid.CancelRequest{
				OrderRef: authResponse.OrderRef,
			})

			require.NoError(t, err)
		})
	}
}

func TestErrorCodes(t *testing.T) {
	ctx := context.Background()

	b, err := bankid.New(bankid.Config{
		URL: bankid.BankIDTestUrl,
		Certificate: bankid.P12Cert{
			Passphrase:    bankid.BankIDTestPassphrase,
			Certificate:   bankid.P12TestCertificate,
			CACertificate: bankid.CATestCertificate,
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
				_, err := b.Auth(ctx, bankid.AuthRequest{
					EndUserIP: "",
				})
				require.Error(t, err)

				_, ok := err.(bankid.RequiredInputMissingError)
				fmt.Printf("type: %T", err)
				require.True(t, ok)

				//using non-existent orderRef
				fake_ref := "non-existent-order-ref"
				_, err = b.Collect(ctx, bankid.CollectRequest{
					OrderRef: fake_ref,
				})
				require.Error(t, err)

				bankIDErr, ok := err.(bankid.BankIDError)
				require.True(t, ok)
				require.Equal(t, bankIDErr.ErrorCode, bankid.InvalidParameters)

			case bankid.AlreadyInProgress:
				ip := randomIPv4()
				personNummer := validPersonnummer()

				//first request
				_, err := b.Auth(ctx, bankid.AuthRequest{
					EndUserIP: ip,
					Requirement: &bankid.Requirement{
						PersonalNumber: personNummer,
					},
				})
				require.NoError(t, err)

				//second request with the same personNummer
				_, err = b.Auth(ctx, bankid.AuthRequest{
					EndUserIP: ip,
					Requirement: &bankid.Requirement{
						PersonalNumber: personNummer,
					},
				})
				require.Error(t, err)
				fmt.Printf("%v", err.Error())

				e, ok := err.(bankid.BankIDError)
				require.True(t, ok)
				require.Equal(t, e.ErrorCode, bankid.AlreadyInProgress)

			case bankid.NotFound:
				b, err := bankid.New(bankid.Config{
					// add non-existent path to URL
					URL: "https://appapi2.test.bankid.com/rp/v6.0/forbidden/path",
					Certificate: bankid.P12Cert{
						Passphrase:    bankid.BankIDTestPassphrase,
						Certificate:   bankid.P12TestCertificate,
						CACertificate: bankid.CATestCertificate,
					},
				})
				require.NoError(t, err)

				_, err = b.Auth(ctx, bankid.AuthRequest{
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

// from: https://github.com/emilybache/personnummer/blob/master/valid_100.txt
func validPersonnummer() string {
	return "3810260632"
}
