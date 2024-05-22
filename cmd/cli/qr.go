package cli

import (
	"fmt"
	"os"

	"github.com/mdp/qrterminal"
	"github.com/nicolaa5/bankid"
	"golang.org/x/term"
)

// Display a QR code in the terminal, generated from QR secret and token received in the /auth and /sign endpoints
func DisplayTerminalQR(qrStartSecret, qrStartToken string, timeDifference int) {
	qrConfig := qrterminal.Config{
		HalfBlocks:     true,
		Level:          qrterminal.L,
		QuietZone:      1,
		Writer:         os.Stdout,
		BlackChar:      qrterminal.WHITE_WHITE,
		BlackWhiteChar: qrterminal.WHITE_BLACK,
		WhiteChar:      qrterminal.BLACK_BLACK,
		WhiteBlackChar: qrterminal.BLACK_WHITE,
	}

	qrCode, err := bankid.GenerateQrPayload(qrStartSecret, qrStartToken, timeDifference)
	if err != nil {
		fmt.Printf("Error generating QR Code: %v\n", err)
		return 
	}

	_, height, err := term.GetSize(0)
	if height < 22 {
		fmt.Print("\033[2J")
		fmt.Printf("\nIncrease the terminal height to see the QR Code\n")
		return
	} else if timeDifference > 0 && height > 22 {
		// removes the height of the QR code (22 lines) from stdout
		fmt.Print("\033[22A\033[J")
	}

	qrterminal.GenerateWithConfig(qrCode, qrConfig)
}
