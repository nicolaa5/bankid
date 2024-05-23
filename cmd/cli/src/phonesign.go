package src

import (
	"fmt"
	"log"

	"github.com/AlecAivazis/survey/v2"
	"github.com/nicolaa5/bankid"
	"github.com/spf13/cobra"
)

var phoneSignRequest bankid.PhoneSignRequest = bankid.PhoneSignRequest{
	Requirement: &bankid.Requirement{},
}

var phoneSignConfig = bankid.Config{
	Certificate: bankid.Certificate{},
}

var phoneSignCommand = &cobra.Command{
	Use:   "phonesign",
	Short: "Sign an agreement over a phone call using BankID",
	Long:  `Use the /phone/sign endpoint to sign a document with BankID from the provided Personnummer`,

	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		if phoneSignConfig.SSLCertificatePath == "" && !test {
			survey.AskOne(
				sslCertPrompt(),
				&phoneSignConfig.SSLCertificatePath,
				survey.WithValidator(survey.Required),
			)
		}

		if phoneSignConfig.Passphrase == "" && !test {
			survey.AskOne(
				passphrasePrompt(),
				&phoneSignConfig.Passphrase,
				survey.WithValidator(survey.Required),
			)
		}

		if phoneSignRequest.PersonalNumber == "" {
			survey.AskOne(
				personalNumberPrompt(),
				&phoneSignRequest.PersonalNumber,
				survey.WithValidator(survey.Required),
			)
		}

		if phoneSignRequest.CallInitiator == "" {
			survey.AskOne(
				callInitiatorPrompt(),
				&phoneSignRequest.CallInitiator,
				survey.WithValidator(survey.Required),
			)
		}

		if phoneSignRequest.UserVisibleData == "" {
			survey.AskOne(
				userDataVisibleDataPrompt(),
				&phoneSignRequest.UserVisibleData,
				survey.WithValidator(survey.Required),
			)
		}

		var b bankid.BankID
		var err error

		if test {
			b, err = bankid.NewTestDefault()
		} else {
			b, err = bankid.New(phoneSignConfig)
		}
		if err != nil {
			log.Fatalf("Internal error in CLI app: %v", err)
		}

		phoneSignResponse, err := b.PhoneSign(ctx, phoneSignRequest)
		if err != nil {
			log.Fatalf("PhoneSign response error: : %v", err)
		}

		fmt.Printf("\n\nWaiting for user to sign a document over phonecall using BankID...\n\n")

		req := bankid.CollectRequest{
			OrderRef: phoneSignResponse.OrderRef,
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
	rootCmd.AddCommand(phoneSignCommand)
	phoneSignCommand.PersistentFlags().BoolVarP(&test, "test", "t", false, "Whether to use the BankID Test environment for the request")
	phoneSignCommand.PersistentFlags().StringVarP(&phoneSignConfig.SSLCertificatePath, "certificatepath", "c", "", "The path to your SSLCertificate")
	phoneSignCommand.PersistentFlags().StringVarP(&phoneSignConfig.Certificate.Passphrase, "passphrase", "p", "", "The password that's used to decrypt the private key of your SSLCertificate")
	phoneSignCommand.PersistentFlags().StringVarP(&phoneSignRequest.CallInitiator, "callinitiator", "i", "", "The user or service that initiated the call")
	phoneSignCommand.PersistentFlags().StringVarP(&phoneSignRequest.PersonalNumber, "personalnumber", "n", "", "The personal number of the individual that will sign the document")
	phoneSignCommand.PersistentFlags().StringVarP(&phoneSignRequest.UserVisibleData, "uservisibledata", "v", "", "The text shown to the enduser")
	phoneSignCommand.PersistentFlags().StringVarP(&phoneSignRequest.UserNonVisibleData, "usernonvisibledata", "w", "", "The provided text is not shown to the enduser")
	phoneSignCommand.PersistentFlags().StringVarP(&phoneSignRequest.UserVisibleDataFormat, "uservisibledataformat", "f", "simpleMarkdownV1", "Format the text that is shown to the user")
	phoneSignCommand.PersistentFlags().StringVarP(&phoneSignRequest.Requirement.CardReader, "cardreader", "e", "", "The transaction must be performed using a card reader where the PIN code is entered on a computer keyboard, or a card reader of higher class")
	phoneSignCommand.PersistentFlags().StringArrayP("certificatepolicies", "s", phoneSignRequest.Requirement.CertificatePolicies, "The oid in certificate policies in the user certificate")
	phoneSignCommand.PersistentFlags().BoolVarP(&phoneSignRequest.Requirement.MRTD, "requiremrtd", "m", false, "Users are required to have a NFC-enabled smartphone which can read MRTD (Machine readable travel document) information to complete the order")
	phoneSignCommand.PersistentFlags().BoolVarP(&phoneSignRequest.Requirement.Pincode, "requirepincode", "o", false, "Users are required to sign the transaction with their PIN code, even if they have biometrics activated.")
	phoneSignCommand.PersistentFlags().StringVarP(&phoneSignRequest.Requirement.PersonalNumber, "requiredpersonalnumber", "u", "", " A personal identification number that is required to complete the transaction")
}
