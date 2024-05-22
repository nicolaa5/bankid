package cli

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/nicolaa5/bankid"
	"github.com/spf13/cobra"
)

var request bankid.AuthRequest = bankid.AuthRequest{
	Requirement: &bankid.Requirement{},
}

var test bool
var certificatepath string
var passphrase string

var authCommand = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate a user with BankID",
	Long:  `Use the /auth endpoint to authenticate a user with BankID and get a QR code to scan.`,

	Run: func(cmd *cobra.Command, args []string) {
		if certificatepath == "" && !test {
			prompt := &survey.Input{
				Message: "Please provide the path to your SSLCertificate",
				Help:    "SSL Certificates help establish trust between your company and BankID by providing a secure connection and validating your organization's identity",
				Suggest: func(toComplete string) []string {
					files, _ := filepath.Glob(toComplete + "*")
					return files
				},
			}

			survey.AskOne(
				prompt,
				&certificatepath,
				survey.WithValidator(survey.Required),
			)
		}

		if passphrase == "" && !test {
			prompt := &survey.Password{
				Message: "Please provide the Passphrase with which your SSLCertificate is encrypted",
				Help:    "The password that is required is used to  protect and secure the private key with which the CSR is signed that BankID uses to generate your SSLCertificate",
			}

			survey.AskOne(
				prompt,
				&passphrase,
				survey.WithValidator(survey.Required),
			)
		}

		if request.EndUserIP == "" {
			prompt := &survey.Input{
				Message: "Please provide the IP address of the Enduser",
				Help:    "EndUserIP is required by BankID because it allows you to improve security by allowing you to compare the enduser's IP you receive with the IP send in the completiondata by BankID.",
				Default: "192.168.0.0",
				Suggest: func(toComplete string) []string {
					return []string{"192.168.0.0", "127.0.0.0"}
				},
			}

			survey.AskOne(
				prompt,
				&request.EndUserIP,
				survey.WithValidator(survey.Required),
			)
		}

		var b bankid.BankID
		if test {
			client, err := bankid.New(bankid.Config{
				URL: bankid.BankIDTestUrl,
				Certificate: bankid.Certificate{
					Passphrase:     bankid.BankIDTestPassphrase,
					SSLCertificate: bankid.SSLTestCertificate,
					CACertificate:  bankid.CATestCertificate,
				},
			})
			if err != nil {
				log.Fatalf("Internal error in CLI app: %v", err)
			}
			b = client

		} else {
			cert, err := bankid.CertificateFromPaths(
				bankid.CertificatePaths{
					Passphrase:         passphrase,
					SSLCertificatePath: certificatepath,
					CACertificatePath:  "../certs/ca_prod.crt",
				},
			)
			if err != nil {
				log.Fatalf("Input error: %v", err)
			}

			client, err := bankid.New(bankid.Config{
				URL:         bankid.BankIDURL,
				Certificate: *cert,
			})
			if err != nil {
				log.Fatalf("Internal error in CLI app: %v", err)
			}
			b = client
		}

		authResponse, err := b.Auth(cmd.Context(), request)
		if err != nil {
			fmt.Printf("Response error: %v\n", err)
			os.Exit(0)
		}

		response := make(chan *bankid.CollectResponse)

		fmt.Printf("\n\nWaiting for user to authenticate using BankID...\n\n")

		req := bankid.CollectRequest{
			OrderRef: authResponse.OrderRef,
		}

		// continuously collect the status of the order
		go b.CollectRoutine(cmd.Context(), req, response)

		start := time.Now().Unix()

		for {
			select {
			case <-cmd.Context().Done():
				return

			case collectResponse, ok := <-response:
				if !ok {
					continue
				}

				now := time.Now().Unix()
				diff := int(now - start)

				displayTerminalQR(authResponse.QrStartSecret, authResponse.QrStartToken, diff)

				if collectResponse.Status == bankid.Complete {
					fmt.Println("Authentication successful")
					prettyPrint(collectResponse.CompletionData.User)
				} else if collectResponse.Status == bankid.Failed {
					fmt.Printf("\nAuthentication failed, reason: %s\n", collectResponse.HintCode)
				}
			default:

			}
		}
	},
}

func init() {
	rootCmd.AddCommand(authCommand)
	authCommand.PersistentFlags().BoolVarP(&test, "test", "t", false, "Whether to use the BankID Test environment for the request")
	authCommand.PersistentFlags().StringVarP(&certificatepath, "certificatepath", "c", "", "The path to your SSLCertificate")
	authCommand.PersistentFlags().StringVarP(&passphrase, "passphrase", "p", "", "The password that's used to decrypt the private key of your SSLCertificate")
	authCommand.PersistentFlags().StringVarP(&request.EndUserIP, "enduserip", "i", "", "The end user's IP address")
	authCommand.PersistentFlags().StringVarP(&request.UserVisibleData, "uservisibledata", "v", "", "The text shown to the enduser")
	authCommand.PersistentFlags().StringVarP(&request.UserNonVisibleData, "usernonvisibledata", "w", "", "The provided text is not shown to the enduser")
	authCommand.PersistentFlags().StringVarP(&request.UserVisibleDataFormat, "uservisibledataformat", "f", "simpleMarkdownV1", "Format the text that is shown to the user")
	authCommand.PersistentFlags().StringVarP(&request.Requirement.CardReader, "cardreader", "e", "", "The transaction must be performed using a card reader where the PIN code is entered on a computer keyboard, or a card reader of higher class")
	authCommand.PersistentFlags().StringArrayP("certificatepolicies", "s", request.Requirement.CertificatePolicies, "The oid in certificate policies in the user certificate")
	authCommand.PersistentFlags().BoolVarP(&request.Requirement.MRTD, "requiremrtd", "m", false, "Users are required to have a NFC-enabled smartphone which can read MRTD (Machine readable travel document) information to complete the order")
	authCommand.PersistentFlags().BoolVarP(&request.Requirement.Pincode, "requirepincode", "o", false, "Users are required to sign the transaction with their PIN code, even if they have biometrics activated.")
	authCommand.PersistentFlags().StringVarP(&request.Requirement.PersonalNumber, "personalnumber", "n", "", " A personal identification number to be used to complete the transaction")
}

func prettyPrint(data interface{}) {
	bytes, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		fmt.Printf(fmt.Sprintf("\n%#v\n", data))
		return
	}
	fmt.Printf(fmt.Sprintf("\n%v\n", string(bytes)))
}
