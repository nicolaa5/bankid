package response

type SignResponse struct {
	// Used to collect the status of the order.
	OrderRef string `json:"orderRef"`
}