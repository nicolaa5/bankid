package response

// ResponseBody is an interface for all successfull BankID responses.
type ResponseBody interface {
	// Unmarshal parses the JSON-encoded data and stores the result in the value pointed to by v.
	Unmarshal(data []byte) error
}

type ErrorResponseBody struct {
	ErrorCode int    `json:"errorCode"`
	Details   string `json:"details"`
}
