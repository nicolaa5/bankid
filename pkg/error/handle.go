package error

func HandleError(err error) {
	switch e := err.(type) {
		case RPError:
			// Handle RPError
		case BankIDError:
			// Handle BankIDError
		default:
			// Handle unknown error
	}
}