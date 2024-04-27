package response

import "encoding/json"

type PhoneAuthResponse struct {
	// Used to collect the status of the order.
	OrderRef string `json:"orderRef"`
}

func (r PhoneAuthResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &r)
}