package response

type PhoneAuthResponse struct {
	// Used to collect the status of the order.
	OrderRef string `json:"orderRef"`
}