package cli

import (
	"fmt"
	"log"

	"github.com/AlecAivazis/survey/v2"
	"github.com/nicolaa5/bankid"
	"github.com/spf13/cobra"
)

var authRequest bankid.AuthRequest = bankid.AuthRequest{
	Requirement: &bankid.Requirement{},
}

var authConfig = bankid.Config{
	Certificate: bankid.Certificate{},
}

var authCommand = &cobra.Command{
	Use:     "auth",
	Aliases: []string{"authenticate"},
	Short:   "Authenticate a user with BankID",
	Long:    `Use the /auth endpoint to authenticate a user with BankID and get a QR code to scan.`,

	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		if authConfig.SSLCertificatePath == "" && !test {
			survey.AskOne(
				sslCertPrompt(),
				&authConfig.SSLCertificatePath,
				survey.WithValidator(survey.Required),
			)
		}

		if authConfig.Passphrase == "" && !test {
			survey.AskOne(
				passphrasePrompt(),
				&authConfig.Passphrase,
				survey.WithValidator(survey.Required),
			)
		}

		if authRequest.EndUserIP == "" {
			survey.AskOne(
				endUserIpPrompt(),
				&authRequest.EndUserIP,
				survey.WithValidator(survey.Required),
			)
		}

		var b bankid.BankID
		var err error 

		if test {
			b, err = bankid.NewTestDefault()
		} else {
			b, err = bankid.New(authConfig)
		}
		if err != nil {
			log.Fatalf("Internal error in CLI app: %v", err)
		}

		authResponse, err := b.Auth(ctx, authRequest)
		if err != nil {
			log.Fatalf("Auth response error: : %v", err)
		}

		fmt.Printf("\n\nWaiting for user to authenticate using BankID...\n\n")

		req := bankid.CollectRequest{
			OrderRef: authResponse.OrderRef,
		}

		response := make(chan *bankid.CollectResponse)

		// continuously collect the status of the order and generate the QR in the terminal
		go b.CollectRoutine(ctx, req, response)
		animateTerminalQR(ctx, authResponse.QrStartSecret, authResponse.QrStartToken, response)
	},
}

func init() {
	rootCmd.AddCommand(authCommand)
	authCommand.PersistentFlags().BoolVarP(&test, "test", "t", false, "Whether to use the BankID Test environment for the request")
	authCommand.PersistentFlags().StringVarP(&authConfig.SSLCertificatePath, "certificatepath", "c", "", "The path to your SSLCertificate")
	authCommand.PersistentFlags().StringVarP(&authConfig.Certificate.Passphrase, "passphrase", "p", "", "The password that's used to decrypt the private key of your SSLCertificate")
	authCommand.PersistentFlags().StringVarP(&authRequest.EndUserIP, "enduserip", "i", "", "The end user's IP address")
	authCommand.PersistentFlags().StringVarP(&authRequest.UserVisibleData, "uservisibledata", "v", "", "The text shown to the enduser")
	authCommand.PersistentFlags().StringVarP(&authRequest.UserNonVisibleData, "usernonvisibledata", "w", "", "The provided text is not shown to the enduser")
	authCommand.PersistentFlags().StringVarP(&authRequest.UserVisibleDataFormat, "uservisibledataformat", "f", "simpleMarkdownV1", "Format the text that is shown to the user")
	authCommand.PersistentFlags().StringVarP(&authRequest.Requirement.CardReader, "cardreader", "e", "", "The transaction must be performed using a card reader where the PIN code is entered on a computer keyboard, or a card reader of higher class")
	authCommand.PersistentFlags().StringArrayP("certificatepolicies", "s", authRequest.Requirement.CertificatePolicies, "The oid in certificate policies in the user certificate")
	authCommand.PersistentFlags().BoolVarP(&authRequest.Requirement.MRTD, "requiremrtd", "m", false, "Users are required to have a NFC-enabled smartphone which can read MRTD (Machine readable travel document) information to complete the order")
	authCommand.PersistentFlags().BoolVarP(&authRequest.Requirement.Pincode, "requirepincode", "o", false, "Users are required to sign the transaction with their PIN code, even if they have biometrics activated.")
	authCommand.PersistentFlags().StringVarP(&authRequest.Requirement.PersonalNumber, "requiredpersonalnumber", "u", "", " A personal identification number that is required to complete the transaction")
}
