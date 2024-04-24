package req

import (
	"bytes"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"time"

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

func (*Client) New(config cfg.Config) (*Client, error) {

	if config.SSLCertificate == nil && config.SSLCertificatePath != "" {
		p12, err := os.ReadFile(config.SSLCertificatePath)
		if err != nil {
			return nil, fmt.Errorf("error reading .p12 file: %w", err)
		}

		config.SSLCertificate = p12
	}

	if config.CACertificate == nil && config.CACertificatePath != "" {
		ca, err := os.ReadFile(config.CACertificatePath)
		if err != nil {
			return nil, fmt.Errorf("error reading root certificate file: %w", err)
		}

		config.CACertificate = ca
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

func (r *Client) req(path string, body []byte) (*http.Response, error) {
    return r.client.Post(
        fmt.Sprintf("%s/%s", r.urlBase, path), 
        "application/json", 
        bytes.NewBuffer(body),
    )
}
