# Start from a base Golang image
FROM golang:1.21-alpine AS builder

LABEL maintainer="Nicolaas Bijman <niekbijman@gmail.com>"
WORKDIR /cli

ENV GOOS=linux 
ENV GOARCH=amd64
ENV CGO_ENABLED=0

# Copy the source code into the container
COPY . /cli

# Build the Golang CLI app
RUN go build -o bankid ./cmd/cli
RUN chmod +x ./bankid

# Final image
FROM alpine
RUN apk add --no-cache bash
RUN apk --no-cache add curl

COPY --from=builder /cli /cli/
WORKDIR /cli

ENTRYPOINT ["./bankid"]