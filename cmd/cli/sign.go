package cli

import (
	"fmt"
	"log"

	"github.com/AlecAivazis/survey/v2"
	"github.com/nicolaa5/bankid"
	"github.com/spf13/cobra"
)

var signRequest bankid.SignRequest = bankid.SignRequest{
	Requirement: &bankid.Requirement{},
}

var signConfig = bankid.Config{
	Certificate: bankid.Certificate{},
}

var signCommand = &cobra.Command{
	Use:   "sign",
	Short: "Sign an agreement using BankID",
	Long:  `Use the /sign endpoint to sign a document or agreement with BankID and get a QR code to scan.`,

	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		if signConfig.SSLCertificatePath == "" && !test {
			survey.AskOne(
				sslCertPrompt(),
				&signConfig.SSLCertificatePath,
				survey.WithValidator(survey.Required),
			)
		}

		if signConfig.Passphrase == "" && !test {
			survey.AskOne(
				passphrasePrompt(),
				&signConfig.Passphrase,
				survey.WithValidator(survey.Required),
			)
		}

		if signRequest.EndUserIP == "" {
			survey.AskOne(
				endUserIpPrompt(),
				&signRequest.EndUserIP,
				survey.WithValidator(survey.Required),
			)
		}

		if signRequest.UserVisibleData == "" {
			survey.AskOne(
				userDataVisibleDataPrompt(),
				&signRequest.UserVisibleData,
				survey.WithValidator(survey.Required),
			)
		}

		var b bankid.BankID
		var err error

		if test {
			b, err = bankid.NewTestDefault()
		} else {
			b, err = bankid.New(signConfig)
		}
		if err != nil {
			log.Fatalf("Internal error in CLI app: %v", err)
		}

		signResponse, err := b.Sign(ctx, signRequest)
		if err != nil {
			log.Fatalf("Sign response error: : %v", err)
		}

		fmt.Printf("\n\nWaiting for user to sign using BankID...\n\n")

		req := bankid.CollectRequest{
			OrderRef: signResponse.OrderRef,
		}

		response := make(chan *bankid.CollectResponse)

		// continuously collect the status of the order and generate the QR in the terminal
		go b.CollectRoutine(ctx, req, response)
		animateTerminalQR(ctx, signResponse.QrStartSecret, signResponse.QrStartToken, response)
	},
}

func init() {
	rootCmd.AddCommand(signCommand)
	signCommand.PersistentFlags().BoolVarP(&test, "test", "t", false, "Whether to use the BankID Test environment for the request")
	signCommand.PersistentFlags().StringVarP(&signConfig.SSLCertificatePath, "certificatepath", "c", "", "The path to your SSLCertificate")
	signCommand.PersistentFlags().StringVarP(&signConfig.Certificate.Passphrase, "passphrase", "p", "", "The password that's used to decrypt the private key of your SSLCertificate")
	signCommand.PersistentFlags().StringVarP(&signRequest.EndUserIP, "enduserip", "i", "", "The end user's IP address")
	signCommand.PersistentFlags().StringVarP(&signRequest.UserVisibleData, "uservisibledata", "v", "", "The text shown to the enduser")
	signCommand.PersistentFlags().StringVarP(&signRequest.UserNonVisibleData, "usernonvisibledata", "w", "", "The provided text is not shown to the enduser")
	signCommand.PersistentFlags().StringVarP(&signRequest.UserVisibleDataFormat, "uservisibledataformat", "f", "simpleMarkdownV1", "Format the text that is shown to the user")
	signCommand.PersistentFlags().StringVarP(&signRequest.Requirement.CardReader, "cardreader", "e", "", "The transaction must be performed using a card reader where the PIN code is entered on a computer keyboard, or a card reader of higher class")
	signCommand.PersistentFlags().StringArrayP("certificatepolicies", "s", signRequest.Requirement.CertificatePolicies, "The oid in certificate policies in the user certificate")
	signCommand.PersistentFlags().BoolVarP(&signRequest.Requirement.MRTD, "requiremrtd", "m", false, "Users are required to have a NFC-enabled smartphone which can read MRTD (Machine readable travel document) information to complete the order")
	signCommand.PersistentFlags().BoolVarP(&signRequest.Requirement.Pincode, "requirepincode", "o", false, "Users are required to sign the transaction with their PIN code, even if they have biometrics activated.")
	signCommand.PersistentFlags().StringVarP(&signRequest.Requirement.PersonalNumber, "requiredpersonalnumber", "u", "", " A personal identification number that is required to complete the transaction")
}
