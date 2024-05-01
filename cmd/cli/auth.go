package cli

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/mdp/qrterminal"
	"github.com/nicolaa5/bankid"
	"github.com/spf13/cobra"
)


var personNummer string
var endUserIp string

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

		response, err := b.Auth(bankid.AuthRequest{
			EndUserIP: endUserIp,
			Requirement: bankid.Requirement{
				PersonalNumber: personNummer,
			},
		})
		if err != nil {
			fmt.Printf("Response error: %v\n", err)
			os.Exit(0)
		}

		prettyPrint(response)

		qrCode, err := b.GenerateQRCode(response.QrStartSecret, response.QrStartToken, 0)
		if err != nil {
			fmt.Printf("Error generating QR Code: %v\n", err)
			os.Exit(0)
		}
		
		config := qrterminal.Config{
			HalfBlocks: true,
			Level: qrterminal.L,
			QuietZone: 2,
			Writer: os.Stdout,
			BlackChar: qrterminal.WHITE_WHITE,
			BlackWhiteChar: qrterminal.WHITE_BLACK,
			WhiteChar: qrterminal.BLACK_BLACK,
			WhiteBlackChar: qrterminal.BLACK_WHITE,
		}
		
		qrterminal.GenerateWithConfig(qrCode, config)
	},
}


func init() {
	rootCmd.AddCommand(authCommand)
	authCommand.PersistentFlags().StringVarP(&personNummer, "personnummer", "p", "190012128888", "The personnummer of the user to authenticate")
	authCommand.PersistentFlags().StringVarP(&endUserIp, "enduserip", "i", "192.168.12.0", "The end user's IP address")
}


func prettyPrint(data interface{}) {
	bytes, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		fmt.Printf(fmt.Sprintf("\n%#v\n", data))
		return 
	}
	fmt.Printf(fmt.Sprintf("\n%v\n", string(bytes)))
}