package request

import "encoding/json"

// Cancels an ongoing sign or auth order. This is typically used if the user cancels the order in your service or app.
type CancelRequest struct {
	// The orderRef returned from auth or sign.
	OrderRef string `json:"orderRef"`
}

func (r CancelRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
