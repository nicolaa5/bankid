package cfg

type Config struct {
	Cert

	// The URL to BankID API
	URL string `json:"url"`

	// The timeout for the request to BankID API in seconds.
	Timeout int `json:"timeout"`
}
