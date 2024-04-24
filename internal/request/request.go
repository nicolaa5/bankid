package request

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/nicolaa5/bankid/pkg/cfg"
)

// BankID is a go client for the BankID API. The RP interface is JSON based.
// 	- HTTP1.1 is required.
// 	- All methods are accessed using HTTP POST.
// 	- HTTP header 'Content-Type' must be set to 'application/json'.
// 	- The parameters including the leading and ending curly bracket is in the body.
type BankIDRequest struct {

}

func New(config cfg.Config) (*BankIDRequest, error) {
    url := "https://www.example.com/api"
    method := "POST"

    // Create a new request
    req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte("param1=value1&param2=value2")))
    if err != nil {
        return nil, fmt.Errorf("error creating new http request: %w", err)
    }

    // Add headers
    req.Header.Set("Content-Type", "application/json")

    // Send the request
    client := &http.Client{}

}