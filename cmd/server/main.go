package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/nicolaa5/bankid"
)

var(
	//go:embed ssl_prod.p12
	SSLProdCertificate []byte
)

func main() {
    e := echo.New()

    err := godotenv.Load(".env")
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    passphrase := os.Getenv("SSL_CERT_PASSPHRASE")
    sslCertPath := os.Getenv("SSL_CERT_PATH")

	b, err := bankid.New(bankid.Config{
		Certificate: bankid.Certificate{
			Passphrase: passphrase,
			SSLCertificatePath: sslCertPath,
		},
	})
	if err != nil {
		log.Fatal("Internal error in new config: %w", err)
	}

    e.GET("/sse/auth", func(c echo.Context) error {
		ctx := c.Request().Context()
        res := c.Response()
        res.Header().Set(echo.HeaderContentType, "text/event-stream")
        res.Header().Set(echo.HeaderCacheControl, "no-cache")
        res.Header().Set(echo.HeaderConnection, "keep-alive")

		authResponse, err := b.Auth(ctx, bankid.AuthRequest{

		})
		if err != nil {
			return fmt.Errorf("auth error: %w", err)
		}

		response := make(chan *bankid.CollectResponse)

		go b.CollectRoutine(
			ctx, 
			bankid.CollectRequest{OrderRef: authResponse.OrderRef}, 
			response,
		)

		for {
			select {
			case <-ctx.Done():
				return nil
	
			case collectResponse, ok := <-response:
				if !ok {
					continue
				}
	
				fmt.Fprintf(res, "data: %#v\n\n", collectResponse)
				res.Flush()
			}
		}
    })

    e.Logger.Fatal(e.Start(":8080"))
}

// func collect(b bankid.BankID, c echo.Context, res *echo.Response, orderRef string) error {

// }