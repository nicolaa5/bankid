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

var test bool
var certificatepath string
var passphrase string
var endUserIp string
var userVisibleData string
var userNonVisibleData string
var userVisibleDataFormat string
var cardReader string
var certificatePolicies []string
var requireMRTD bool
var requirePinCode bool
var personalNumber string

var authCommand = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate a user with BankID",
	Long:  `Use the /auth endpoint to authenticate a user with BankID and get a QR code to scan.`,

	Run: func(cmd *cobra.Command, args []string) {
		if certificatepath == ""  && !test {
			prompt := &survey.Input{
				Message: "Please provide the path to your SSLCertificate",
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
			}

			survey.AskOne(
				prompt,
				&passphrase,
				survey.WithValidator(survey.Required),
			)
		}

		if endUserIp == "" {
			prompt := &survey.Input{
				Message: "Please provide the IP address of the Enduser",
				Help: "EndUserIP is required by BankID because it allows you to improve security by allowing you to compare the enduser's IP you receive with the IP send in the completiondata by BankID.",
				Default: "192.168.0.0",
				Suggest: func(toComplete string) []string {
					return []string{"192.168.0.0", "127.0.0.0"}
				},
			}

			survey.AskOne(
				prompt,
				&endUserIp,
				survey.WithValidator(survey.Required),
			)
		}

		var b bankid.BankID
		if test {
			client, err := bankid.New(bankid.Parameters{
				URL:         bankid.BankIDTestUrl,
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

			client, err := bankid.New(bankid.Parameters{
				URL:         bankid.BankIDURL,
				Certificate: *cert,
			})
			if err != nil {
				log.Fatalf("Internal error in CLI app: %v", err)
			}
			b = client
		}

		request, err := bankid.NewRequest[bankid.AuthRequest](
			bankid.WithEndUserIP(endUserIp),
			bankid.WithUserVisibleData(userVisibleData),
			bankid.WithUserNonVisibleData(userNonVisibleData),
			bankid.WithUserVisibleDataFormat(userVisibleDataFormat),
			bankid.WithCardReader(cardReader),
			bankid.WithCertificatePolicies(certificatePolicies),
			bankid.WithMRTD(requireMRTD),
			bankid.WithPersonalNumber(personalNumber),
			bankid.WithPincode(requirePinCode),
		)
		if err != nil {
			log.Fatalf("New Request: %v", err.Error())
		}

		authResponse, err := b.Auth(request)
		if err != nil {
			fmt.Printf("Response error: %v\n", err)
			os.Exit(0)
		}

		response := make(chan *bankid.CollectResponse)
		quit := make(chan struct{})

		// keep collecting the status of the order
		go CollectRoutine(b, authResponse.OrderRef, response, quit)

		start := time.Now().Unix()
		fmt.Printf("\n\nWaiting for user to authenticate using BankID...\n\n")
		
		for {
			select {
			case collectResponse := <-response:
				now := time.Now().Unix()
				diff := int(now - start)
		
				DisplayTerminalQR(authResponse.QrStartSecret, authResponse.QrStartToken, diff)

				if collectResponse.Status == bankid.Complete {
					fmt.Println("Authentication successful")
					prettyPrint(collectResponse.CompletionData)
				} else if collectResponse.Status == bankid.Failed {					
					fmt.Printf("\nAuthentication failed, reason: %s\n", collectResponse.HintCode)
				}
				
			case <-quit:
				return
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
	authCommand.PersistentFlags().StringVarP(&endUserIp, "enduserip", "i", "", "The end user's IP address")
	authCommand.PersistentFlags().StringVarP(&userVisibleData, "uservisibledata", "v", "", "The text shown to the enduser")
	authCommand.PersistentFlags().StringVarP(&userNonVisibleData, "usernonvisibledata", "w", "", "The provided text is not shown to the enduser")
	authCommand.PersistentFlags().StringVarP(&userVisibleDataFormat, "uservisibledataformat", "f", "simpleMarkdownV1", "Format the text that is shown to the user")
	authCommand.PersistentFlags().StringVarP(&cardReader, "cardreader", "e", "", "The transaction must be performed using a card reader where the PIN code is entered on a computer keyboard, or a card reader of higher class")
	authCommand.PersistentFlags().StringArrayP("certificatepolicies", "s", certificatePolicies, "The oid in certificate policies in the user certificate")
	authCommand.PersistentFlags().BoolVarP(&requireMRTD, "requiremrtd", "m", false, "Users are required to have a NFC-enabled smartphone which can read MRTD (Machine readable travel document) information to complete the order")
	authCommand.PersistentFlags().BoolVarP(&requirePinCode, "requirepincode", "o", false, "Users are required to sign the transaction with their PIN code, even if they have biometrics activated.")
	authCommand.PersistentFlags().StringVarP(&personalNumber, "personalnumber", "n", "", " A personal identification number to be used to complete the transaction")
}

func prettyPrint(data interface{}) {
	bytes, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		fmt.Printf(fmt.Sprintf("\n%#v\n", data))
		return
	}
	fmt.Printf(fmt.Sprintf("\n%v\n", string(bytes)))
}
