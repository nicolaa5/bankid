package bankid

import (
	"bytes"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"software.sslmate.com/src/go-pkcs12"
)

// BankID is a go client for the BankID API. The RP interface is JSON based.
//   - HTTP1.1 is required.
//   - All methods are accessed using HTTP POST.
//   - HTTP header 'Content-Type' must be set to 'application/json'.
//   - The parameters including the leading and ending curly bracket is in the body.
type Config struct {
	UrlBase string
	Client  *http.Client
}

type ClientParameters struct {
	Path   string
	Config *Config
	Body   RequestBody
}

// request sends a request to the BankID API and returns the response.
func request[T ResponseBody](p ClientParameters) (r *T, err error) {
	b, err := p.Body.Marshal()
	if err != nil {
		return nil, fmt.Errorf("error marshalling body: %w", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", p.Config.UrlBase, p.Path), bytes.NewBuffer(b))
	if err != nil {

		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Host", "appapi2.bankid.com")
	req.Header.Set("Content-Type", "application/json")

	res, err := p.Config.Client.Do(req)
	if res == nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	switch err.(type) {
	case nil:
		// All good
	default:
		e := BankIDError{}
		err := json.Unmarshal(body, &e)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling error response: %w", err)
		}

		return nil, assignError(e.ErrorCode)
	}

	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	return r, nil
}

func newConfig(config Parameters) (*Config, error) {
	// Decode the .p12 certificate into a private key and the certificate
	key, cert, err := pkcs12.Decode(config.SSLCertificate, config.Passphrase)
	if err != nil {
		return nil, fmt.Errorf("error decoding certificate: %w", err)
	}

	privateKey, ok := key.(crypto.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("could not assert private key")
	}

	certPool := x509.NewCertPool()
	ok = certPool.AppendCertsFromPEM(config.CACertificate)
	if !ok {
		return nil, fmt.Errorf("could not append root certificate to pool")
	}

	// Create a new TLS configuration with the key and certificate
	tlsConfig := &tls.Config{
		RootCAs: certPool,
		Certificates: []tls.Certificate{
			{
				Certificate: [][]byte{cert.Raw},
				PrivateKey:  privateKey,
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

	return &Config{
		UrlBase: config.URL,
		Client:  client,
	}, nil
}
