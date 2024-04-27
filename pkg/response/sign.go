package response

import "encoding/json"

type SignResponse struct {
	// Used to collect the status of the order.
	OrderRef string `json:"orderRef"`
}

func (r SignResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &r)
}