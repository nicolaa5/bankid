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
	authCommand.PersistentFlags().StringVarP(&userVisibleData, "uservisibledata", "u", "", "The text shown to the enduser")
}

func prettyPrint(data interface{}) {
	bytes, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		fmt.Printf(fmt.Sprintf("\n%#v\n", data))
		return 
	}
	fmt.Printf(fmt.Sprintf("\n%v\n", string(bytes)))
}