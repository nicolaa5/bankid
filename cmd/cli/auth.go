package cli

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	// "github.com/AlecAivazis/survey/v2"
	"github.com/mdp/qrterminal"
	"github.com/nicolaa5/bankid"
	"github.com/spf13/cobra"
)


var endUserIp string
var personNummer string
var userVisibleData string
var userNonVisibleData string
var userVisibleDataFormat string
var cardReader string
var certificatePolicies []string
var requireMRTD bool
var requirePinCode bool
var personalNumber string
var risk string
var returnRisk string
var returnUrl string

var authCommand = &cobra.Command{
	Use:   "auth",
	Short: "Call the /auth endpoint to authenticate a user",
	Long: `Use the /auth endpoint to authenticate a user and get a QR code to scan.`,

	Run: func(cmd *cobra.Command, args []string) {
		b, err := bankid.New(bankid.Parameters{
			URL: bankid.BankIDTestUrl,
			Certificate: bankid.Certificate{
				Passphrase: bankid.BankIDTestPassphrase,
				SSLCertificate: bankid.SSLTestCertificate,
				CACertificate: bankid.CATestCertificate,
			},
		})
		if err != nil {
			log.Fatalf("Internal error in CLI app: %v", err)
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
			bankid.WithReturnRisk(returnRisk),
			bankid.WithReturnUrl(returnUrl),
			bankid.WithRisk(risk),
		)
		if err != nil {
			log.Fatalf("%v", err.Error())
		}

		authResponse, err := b.Auth(request)
		if err != nil {
			fmt.Printf("Response error: %v\n", err)
			os.Exit(0)
		}
			
		qrConfig := qrterminal.Config{
			HalfBlocks: true,
			Level: qrterminal.L,
			QuietZone: 1,
			Writer: os.Stdout,
			BlackChar: qrterminal.WHITE_WHITE,
			BlackWhiteChar: qrterminal.WHITE_BLACK,
			WhiteChar: qrterminal.BLACK_BLACK,
			WhiteBlackChar: qrterminal.BLACK_WHITE,
		}

		fmt.Printf("\n\nWaiting for authentication...\n\n")

		start := time.Now().Unix()

		for {
			now := time.Now().Unix()
			diff := int(now - start)

			qrCode, err := b.GenerateQRCode(authResponse.QrStartSecret, authResponse.QrStartToken, diff)
			if err != nil {
				fmt.Printf("Error generating QR Code: %v\n", err)
				os.Exit(0)
			}

			if diff != 0 {
				// removes the lines of the last QR code from stdout
				fmt.Print("\033[22A\033[J")
			}

			qrterminal.GenerateWithConfig(qrCode, qrConfig)

			collectResponse, err := b.Collect(bankid.CollectRequest{
				OrderRef: authResponse.OrderRef,
			})
			if err != nil {
				fmt.Printf("Error collecting status: %v\n", err)
				os.Exit(0)
			}

			if collectResponse.Status == bankid.Pending {
				time.Sleep(1 * time.Second)
				continue
			}

			if collectResponse.Status == bankid.Complete {
				fmt.Println("Authentication successful")

				prettyPrint(collectResponse)
				break
			}

			fmt.Println("Authentication failed")
			prettyPrint(collectResponse)
			break
		}
	},
}


func init() {
	rootCmd.AddCommand(authCommand)
	authCommand.PersistentFlags().StringVarP(&personNummer, "personnummer", "p", "", "The personnummer of the user to authenticate")
	authCommand.PersistentFlags().StringVarP(&endUserIp, "enduserip", "i", "", "The end user's IP address")
	authCommand.PersistentFlags().StringVarP(&userVisibleData, "uservisibledata", "v", "", "The text shown to the enduser")
	authCommand.PersistentFlags().StringVarP(&userNonVisibleData, "usernonvisibledata", "w", "", "The provided text is not shown to the enduser")
	authCommand.PersistentFlags().StringVarP(&userVisibleDataFormat, "uservisibledataformat", "f", "simpleMarkdownV1", "Format the text that is shown to the user")
	authCommand.PersistentFlags().StringVarP(&cardReader, "cardreader", "e", "", "The transaction must be performed using a card reader where the PIN code is entered on a computer keyboard, or a card reader of higher class")
	authCommand.PersistentFlags().StringArrayP("certificatepolicies", "c", certificatePolicies, "The oid in certificate policies in the user certificate")
	authCommand.PersistentFlags().BoolVarP(&requireMRTD, "requiremrtd", "m", false, "the client needs to provide MRTD (Machine readable travel document) information to complete the order")
	authCommand.PersistentFlags().BoolVarP(&requirePinCode, "requirepincode", "o", false, "Users are required to sign the transaction with their PIN code, even if they have biometrics activated.")
	authCommand.PersistentFlags().StringVarP(&personalNumber, "personalnumber", "n", "", " A personal identification number to be used to complete the transaction")
	authCommand.PersistentFlags().StringVarP(&risk, "risk", "r", "", "Set the acceptable risk level for the transaction.")
	authCommand.PersistentFlags().StringVarP(&returnRisk, "returnrisk", "x", "", "a risk indication will be included in the collect response when the order completes")
	authCommand.PersistentFlags().StringVarP(&returnUrl, "returnurl", "u", "", "Orders started on the same device (started with autostart token) will call this URL when the order is completed")
}

func prettyPrint(data interface{}) {
	bytes, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		fmt.Printf(fmt.Sprintf("\n%#v\n", data))
		return 
	}
	fmt.Printf(fmt.Sprintf("\n%v\n", string(bytes)))
}