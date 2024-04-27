package response

import "encoding/json"

type PhoneSignResponse struct {
	// Used to collect the status of the order.
	OrderRef string `json:"orderRef"`
}

func (r PhoneSignResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &r)
}