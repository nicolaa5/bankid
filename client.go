package bankid

import (
	"bytes"
	"context"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/youmark/pkcs8"
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

func newRequestConfig(params Config) (*RequestConfig, error) {
	var cert *tls.Certificate

	switch v := params.Certificate.(type) {
	case P12Cert:
		c, err := decodeP12(v)
		if err != nil {
			return nil, fmt.Errorf("decode P12 error: %w", err)
		}
		cert = c
	case PEMCert:
		c, err := decodePEM(v)
		if err != nil {
			return nil, fmt.Errorf("decode PEM error: %w", err)
		}
		cert = c
	}

	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM(params.CA())
	if !ok {
		return nil, fmt.Errorf("could not append root certificate to pool")
	}

	// Create a new TLS configuration with the key and certificate
	tlsConfig := &tls.Config{
		RootCAs: certPool,
		Certificates: []tls.Certificate{
			*cert,
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

// PEM format for BankID is the .p12 converted to .pem.
func decodePEM(c PEMCert) (*tls.Certificate, error) {
	publicKey, privateKey, err := parsePem([]byte(c.Certificate), c.Passphrase)
	if err != nil {
		return nil, fmt.Errorf("error parsing private key: %w", err)
	}

	cert := &tls.Certificate{
		Certificate: [][]byte{publicKey.Raw},
		PrivateKey:  privateKey,
		Leaf:        publicKey,
	}

	return cert, nil
}

func parsePem(bytes []byte, passphrase string) (*x509.Certificate, crypto.PrivateKey, error) {
	var publicKey *x509.Certificate
	var privateKey interface{}

	for {
		block, rest := pem.Decode(bytes)
		if block == nil {
			break
		}
		bytes = rest

		switch block.Type {
		case "CERTIFICATE":
			x509Cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return nil, nil, fmt.Errorf("error parsing certificate: %w", err)
			}

			publicKey = x509Cert

		case "ENCRYPTED PRIVATE KEY":
			k, err := decryptPrivateKey(block.Bytes, passphrase)
			if err != nil {
				return nil, nil, fmt.Errorf("could not decrypt private key: %w", err)
			}

			privateKey = k

		case "PRIVATE KEY":
			key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to parse PKCS#8 private key: %v", err)
			}
			privateKey = key

		case "EC PRIVATE KEY":
			key, err := x509.ParseECPrivateKey(block.Bytes)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to parse EC private key: %v", err)
			}
			privateKey = key

		case "RSA PRIVATE KEY":
			key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to parse RSA private key: %v", err)
			}
			privateKey = key
		}
	}

	key, ok := privateKey.(crypto.PrivateKey)
	if !ok {
		return nil, nil, fmt.Errorf("could not assert private key")
	}

	return publicKey, key, nil
}

func decryptPrivateKey(pem []byte, passphrase string) (crypto.PrivateKey, error) {
	key, err := pkcs8.ParsePKCS8PrivateKey(pem, []byte(passphrase))
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt PEM block: %v", err)
	}

	privateKey, ok := key.(crypto.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("could not assert private key")
	}

	return privateKey, nil
}

func decodeP12(c P12Cert) (*tls.Certificate, error) {
	key, x509Cert, err := pkcs12.Decode(c.Certificate, c.Passphrase)
	if err != nil {
		return nil, fmt.Errorf("error decoding P12 certificate: %w", err)
	}

	privateKey, ok := key.(crypto.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("could not assert private key")
	}

	cert := tls.Certificate{
		Certificate: [][]byte{x509Cert.Raw},
		PrivateKey:  privateKey,
		Leaf:        x509Cert,
	}

	return &cert, nil
}
