package bankid

import (
	"bytes"
	"context"
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

// RequestConfig contains the required config for each request to BankID
//   - HTTP1.1 is required.
//   - All methods are accessed using HTTP POST.
//   - HTTP header 'Content-Type' must be set to 'application/json'.
type RequestConfig struct {
	UrlBase string
	Client  *http.Client
}

type RequestParameters struct {
	Path   string
	Config *RequestConfig
	Body   RequestBody
}

// request sends a request to the BankID API and handles and returns the response or error.
func request[T ResponseBody](ctx context.Context, p RequestParameters) (r *T, err error) {
	b, err := p.Body.Marshal()
	if err != nil {
		return nil, fmt.Errorf("error marshalling body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s%s", p.Config.UrlBase, p.Path), bytes.NewBuffer(b))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := p.Config.Client.Do(req)
	if res == nil {
		return nil, fmt.Errorf("error request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	if res.StatusCode >= 300 {
		e := BankIDError{}
		err := json.Unmarshal(body, &e)
		if err != nil {
			return nil, fmt.Errorf("unknown error: %w", err)
		}

		return nil, assignError(e.ErrorCode)
	}

	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	return r, nil
}

func newRequestConfig(params Parameters) (*RequestConfig, error) {
	// Decode the .p12 certificate into a private key and the certificate
	key, cert, err := pkcs12.Decode(params.SSLCertificate, params.Passphrase)
	if err != nil {
		return nil, fmt.Errorf("error decoding certificate: %w", err)
	}

	privateKey, ok := key.(crypto.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("could not assert private key")
	}

	certPool := x509.NewCertPool()
	ok = certPool.AppendCertsFromPEM(params.CACertificate)
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
		Timeout: time.Second * time.Duration(params.Timeout),
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	return &RequestConfig{
		UrlBase: params.URL,
		Client:  client,
	}, nil
}
