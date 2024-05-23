package cli

import (
	"fmt"
	"log"

	"github.com/AlecAivazis/survey/v2"
	"github.com/nicolaa5/bankid"
	"github.com/spf13/cobra"
)

var phoneAuthRequest bankid.PhoneAuthRequest = bankid.PhoneAuthRequest{
	Requirement: &bankid.Requirement{},
}

var phoneAuthConfig = bankid.Config{
	Certificate: bankid.Certificate{},
}

var phoneAuthCommand = &cobra.Command{
	Use:   "phoneauth",
	Short: "Authenticate a user over a phone call with BankID",
	Long:  `Use the /phone/auth endpoint to authenticate a user with BankID from the provided Personnummer`,

	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		if phoneAuthConfig.SSLCertificatePath == "" && !test {
			survey.AskOne(
				sslCertPrompt(),
				&phoneAuthConfig.SSLCertificatePath,
				survey.WithValidator(survey.Required),
			)
		}

		if phoneAuthConfig.Passphrase == "" && !test {
			survey.AskOne(
				passphrasePrompt(),
				&phoneAuthConfig.Passphrase,
				survey.WithValidator(survey.Required),
			)
		}

		if phoneAuthRequest.PersonalNumber == "" {
			survey.AskOne(
				personalNumberPrompt(),
				&phoneAuthRequest.PersonalNumber,
				survey.WithValidator(survey.Required),
			)
		}

		if phoneAuthRequest.CallInitiator == "" {
			survey.AskOne(
				callInitiatorPrompt(),
				&phoneAuthRequest.CallInitiator,
				survey.WithValidator(survey.Required),
			)
		}

		var b bankid.BankID
		var err error

		if test {
			b, err = bankid.NewTestDefault()
		} else {
			b, err = bankid.New(phoneAuthConfig)
		}
		if err != nil {
			log.Fatalf("Internal error in CLI app: %v", err)
		}

		phoneAuthResponse, err := b.PhoneAuth(ctx, phoneAuthRequest)
		if err != nil {
			log.Fatalf("PhoneAuth response error: : %v", err)
		}

		fmt.Printf("\n\nWaiting for user to phoneAuthenticate using BankID...\n\n")

		req := bankid.CollectRequest{
			OrderRef: phoneAuthResponse.OrderRef,
		}

		response := make(chan *bankid.CollectResponse)

		// continuously collect the status of the order and generate the QR in the terminal
		go b.CollectRoutine(ctx, req, response)

		for {
			select {
			case <-ctx.Done():
				return

			case collectResponse, ok := <-response:
				if !ok {
					continue
				}

				prettyPrint(collectResponse)
			default:

			}
		}
	},
}

func init() {
	rootCmd.AddCommand(phoneAuthCommand)
	phoneAuthCommand.PersistentFlags().BoolVarP(&test, "test", "t", false, "Whether to use the BankID Test environment for the request")
	phoneAuthCommand.PersistentFlags().StringVarP(&phoneAuthConfig.SSLCertificatePath, "certificatepath", "c", "", "The path to your SSLCertificate")
	phoneAuthCommand.PersistentFlags().StringVarP(&phoneAuthConfig.Certificate.Passphrase, "passphrase", "p", "", "The password that's used to decrypt the private key of your SSLCertificate")
	phoneAuthCommand.PersistentFlags().StringVarP(&phoneAuthRequest.CallInitiator, "callinitiator", "i", "", "The user or service that initiated the call")
	phoneAuthCommand.PersistentFlags().StringVarP(&phoneAuthRequest.PersonalNumber, "personalnumber", "n", "", "The personal number of the individual that will be authenticated")
	phoneAuthCommand.PersistentFlags().StringVarP(&phoneAuthRequest.UserVisibleData, "uservisibledata", "v", "", "The text shown to the enduser")
	phoneAuthCommand.PersistentFlags().StringVarP(&phoneAuthRequest.UserNonVisibleData, "usernonvisibledata", "w", "", "The provided text is not shown to the enduser")
	phoneAuthCommand.PersistentFlags().StringVarP(&phoneAuthRequest.UserVisibleDataFormat, "uservisibledataformat", "f", "simpleMarkdownV1", "Format the text that is shown to the user")
	phoneAuthCommand.PersistentFlags().StringVarP(&phoneAuthRequest.Requirement.CardReader, "cardreader", "e", "", "The transaction must be performed using a card reader where the PIN code is entered on a computer keyboard, or a card reader of higher class")
	phoneAuthCommand.PersistentFlags().StringArrayP("certificatepolicies", "s", phoneAuthRequest.Requirement.CertificatePolicies, "The oid in certificate policies in the user certificate")
	phoneAuthCommand.PersistentFlags().BoolVarP(&phoneAuthRequest.Requirement.MRTD, "requiremrtd", "m", false, "Users are required to have a NFC-enabled smartphone which can read MRTD (Machine readable travel document) information to complete the order")
	phoneAuthCommand.PersistentFlags().BoolVarP(&phoneAuthRequest.Requirement.Pincode, "requirepincode", "o", false, "Users are required to sign the transaction with their PIN code, even if they have biometrics activated.")
	phoneAuthCommand.PersistentFlags().StringVarP(&phoneAuthRequest.Requirement.PersonalNumber, "requiredpersonalnumber", "u", "", " A personal identification number that is required to complete the transaction")
}
