package response

// A successful response contains an empty JSON object.
type CancelResponse struct{}

func (r CancelResponse) Unmarshal(data []byte) error {
	return nil
}
