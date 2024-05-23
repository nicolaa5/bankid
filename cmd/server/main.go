package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/nicolaa5/bankid"
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

    e.POST("/sse/auth", func(c echo.Context) error {
		ctx := c.Request().Context()
        res := c.Response()
		
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return fmt.Errorf("read request body error: %w", err)
		}

		authRequest := bankid.AuthRequest{}

		if err := json.Unmarshal(body, &authRequest); err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("Invalid input error: %s", string(body)))
		}

		authResponse, err := b.Auth(ctx, authRequest)
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Auth error: %v", err.Error()))
		}

		response := make(chan *bankid.CollectResponse)

		go b.CollectRoutine(
			ctx, 
			bankid.CollectRequest{OrderRef: authResponse.OrderRef}, 
			response,
		)

        res.Header().Set(echo.HeaderContentType, "text/event-stream")
        res.Header().Set(echo.HeaderCacheControl, "no-cache")
        res.Header().Set(echo.HeaderConnection, "keep-alive")

		for {
			select {
			case <-ctx.Done():
				return nil
	
			case collectResponse, ok := <-response:
				if !ok {
					continue
				}

				bytes, err := json.Marshal(collectResponse)
				if err != nil {
					c.String(http.StatusInternalServerError, fmt.Sprintf("Data marshal error: %v", err.Error()))
				}
	
				fmt.Fprintf(res, "data: %s\n\n",  string(bytes))
				res.Flush()
			}
		}
    })

    e.Logger.Fatal(e.Start(":8080"))
}