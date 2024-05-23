package src

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var test bool

var rootCmd = &cobra.Command{
	Use:   "bankid",
	Short: "Command line interface to interact with bankid package",
	Long:  `Use this CLI as a frontend client to generate QR codes and interact with the BankID API.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()

		// Define survey questions
		var endpoint string
		prompt := &survey.Select{
			Message: "Please choose the BankID endpoint to call:",
			Options: []string{
				"/auth",
				"/sign",
				"/phone/auth",
				"/phone/sign",
			},
		}
		survey.AskOne(prompt, &endpoint, survey.WithValidator(survey.Required))

		switch endpoint {
		case "/auth":
			authCommand.Run(cmd, args)
		case "/sign":
			signCommand.Run(cmd, args)
		case "/phone/auth":
			phoneAuthCommand.Run(cmd, args)
		case "/phone/sign":
			phoneSignCommand.Run(cmd, args)
		}

		// Process user input
		fmt.Printf("Called, %s endpoint\n", endpoint)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func prettyPrint(data interface{}) {
	bytes, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		fmt.Printf(fmt.Sprintf("\n%#v\n", data))
		return
	}
	fmt.Printf(fmt.Sprintf("\n%v\n", string(bytes)))
}
