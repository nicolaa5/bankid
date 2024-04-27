package client

import (
	"bytes"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/nicolaa5/bankid/internal/validate"
	"github.com/nicolaa5/bankid/pkg/cfg"
	"software.sslmate.com/src/go-pkcs12"
)

// BankID is a go client for the BankID API. The RP interface is JSON based.
//   - HTTP1.1 is required.
//   - All methods are accessed using HTTP POST.
//   - HTTP header 'Content-Type' must be set to 'application/json'.
//   - The parameters including the leading and ending curly bracket is in the body.
type Client struct {
	urlBase string
	client  *http.Client
}

func (r *Client) Request(path string, body []byte) (*http.Response, error) {
	body, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("error marshalling body: %w", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", r.urlBase, path), bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	return r.client.Do(req)
}


func New(config cfg.Config) (*Client, error) {
	if err := validate.Config(&config); err != nil {
		return nil, fmt.Errorf("error validating config: %w", err)
	}

	err := setCert(&config)
	if err != nil {
		return nil, fmt.Errorf("error setting certificate: %w", err)
	}

	// Parse the decrypted .p12 data
	privateKey, cert, err := pkcs12.Decode(config.CACertificate, config.Passphrase)
	if err != nil {
		return nil, fmt.Errorf("error decoding certificate: %w", err)
	}

	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM(config.CACertificate)
	if !ok {
		return nil, fmt.Errorf("could not append root certificate to pool")
	}

	// Create a new TLS configuration with the key and certificate
	tlsConfig := &tls.Config{
		RootCAs: certPool,
		Certificates: []tls.Certificate{
			{
				Certificate: [][]byte{cert.Raw},
				PrivateKey:  privateKey.(crypto.PrivateKey),
				Leaf:        cert,
			},
		},
	}

	// Create an HTTP client with the custom TLS configuration
	client := &http.Client{
		Timeout: time.Second * time.Duration(config.Timeout),
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	return &Client{
		urlBase: config.URL,
		client:  client,
	}, nil
}